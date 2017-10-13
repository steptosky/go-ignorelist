/*
**  Copyright(C) 2017, StepToSky
**
**  Redistribution and use in source and binary forms, with or without
**  modification, are permitted provided that the following conditions are met:
**
**  1.Redistributions of source code must retain the above copyright notice, this
**    list of conditions and the following disclaimer.
**  2.Redistributions in binary form must reproduce the above copyright notice,
**    this list of conditions and the following disclaimer in the documentation
**    and / or other materials provided with the distribution.
**  3.Neither the name of StepToSky nor the names of its contributors
**    may be used to endorse or promote products derived from this software
**    without specific prior written permission.
**
**  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
**  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
**  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
**  DISCLAIMED.IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
**  ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
**  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
**  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
**  ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
**  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
**  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
**
**  Contacts: www.steptosky.com
 */

package ignore

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

const (
	not1                 = "not "
	not2                 = "!"
	pathSeparator string = string(os.PathSeparator)
)

var path_sep_replacer_regex *regexp.Regexp = regexp.MustCompile("[\\\\/:]+")

type pattern struct {
	tag    string
	prefix string
	suffix string
	isFile bool
}

func (s *pattern) HasPrefix() bool {
	return len(s.prefix) != 0
}

func (s *pattern) HasSuffix() bool {
	return len(s.suffix) != 0
}

func (s *pattern) IsEmpty() bool {
	return !s.HasPrefix() && !s.HasSuffix()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

// Ignore list.
// Common pattern conception is "prefix*suffix".
//
// Example of the ignore patterns:
//
// some-folder/*- Ignores folder with name "some-folder" and all its children.
// some-folder/*.ex - Ignores files with extension ".ex" in folder "some-folder" and all its children..
// *.ex - Ignores files with extension ".ex" in all folders.
// some-folder/file - Ignores file with the name "file" in folder "some-folder". I.e. ignoring files by its full path.
//
// For including files and folder you can use "not " or "!":
//
// some-folder/*
// not some-folder/*.ex
//
// The example above ignores all files and folders within "some-folder" except all files with extension .ex,
// it also effects children folders too.
//
// ATTENTION between "not" and next symbols must be one space "not some-folder",
// the "!" does not need that space "!some-folder".
//
// The path separator can be one of the following symbols: \ / :
//
// The tag usage example:
// [Any text] some-folder/*.ex
// You can get the tag with method IsIgnoredEx
// You can use the tags it as you wish for any porpoises.
// The ignore list does not use tags at all, it just extract it for you.

type List struct {
	excludePatternList []pattern
	includePatternList []pattern
}

// Returns new ignore list.
func NewList() *List {
	return &List{}
}

// It returns new ignore list even if an error is occurred.
func NewListFromFile(filePath string) (*List, error) {
	list := NewList()
	if err := list.LoadFromFile(filePath); err != nil {
		return list, err
	}
	return list, nil
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

// Loads ignore list data from specified file.
// It clears the struct before loading the file.
// Use combine method if you want to combine 2 ignore lists.
// It returns an error if occurred that describes the problem and makes the struct empty.
func (ignoreList *List) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	ignoreList.Clear()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err = ignoreList.processLine(&line); err != nil {
			ignoreList.Clear()
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Combines 2 ignore lists in one this.
// If this ignore list already contains th same pattern from other list
// then pattern from other list will be used and will replace "tag".
func (ignoreList *List) Combine(otherIgnoreList *List) *List {
	ignoreList.combine(&ignoreList.includePatternList, &otherIgnoreList.includePatternList)
	ignoreList.combine(&ignoreList.excludePatternList, &otherIgnoreList.excludePatternList)
	return ignoreList
}

// Adds a new pattern to the ignore list.
// It does not check if the same patter is already exist.
func (ignoreList *List) AddPattern(pattern string) error {
	if err := ignoreList.processLine(&pattern); err != nil {
		return err
	}
	return nil
}

// It returns true if the given file path in the ignore list otherwise false.
// See IsIgnoredEx
func (ignoreList *List) IsIgnored(filePath string) bool {
	res, _ := ignoreList.IsIgnoredEx(filePath)
	return res
}

// It returns true if the given file path in the ignore list otherwise false.
// Also it returns a tag. The tag depends on your rules, you can use it as you wish.
//
// Note: include list is the first processed,
// so if the file is found in the include list the algorithm will not try to find it in the exclude list.
// I.e. if there are 2 patterns like this: "not folder/1" and "folder/1" the result
// "folder/1" is always included (i.e. the function returns false)
func (ignoreList *List) IsIgnoredEx(filePath string) (bool, string) {
	if len(ignoreList.includePatternList) == 0 && len(ignoreList.excludePatternList) == 0 {
		return false, ""
	}
	//------------
	res, idx := ignoreList.hasMatchedPattern(&filePath, &ignoreList.includePatternList)
	if res {
		return false, ignoreList.includePatternList[idx].tag
	}
	//------------
	res, idx = ignoreList.hasMatchedPattern(&filePath, &ignoreList.excludePatternList)
	if res {
		return true, ignoreList.excludePatternList[idx].tag
	}
	//------------
	return false, ""
}

// Clears ignore list.
func (ignoreList *List) Clear() {
	if len(ignoreList.excludePatternList) != 0 {
		ignoreList.excludePatternList = ignoreList.excludePatternList[:0]
	}
	if len(ignoreList.includePatternList) != 0 {
		ignoreList.includePatternList = ignoreList.includePatternList[:0]
	}
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func extractTag(str *string) (string, string, error) {
	outLine := ""
	outTag := ""
	strLen := len(*str)
	if strLen == 0 {
		return outLine, outTag, nil
	}
	if strings.HasPrefix(*str, "[") {
		idx := strings.Index(*str, "]")
		if idx == -1 {
			return outLine, outTag, errors.New("tag is not closed, you must use <]> symbol to close it")
		}
		if idx != strLen-1 {
			outLine = (*str)[idx+1: strLen]
		}
		outTag = (*str)[1:idx]
	} else {
		outLine = *str
	}
	return outLine, outTag, nil
}

func removeNot(str *string) *string {
	outStr := strings.TrimPrefix(*str, not1)
	outStr = strings.TrimPrefix(outStr, not2)
	return &outStr
}

func fixSeparator(str *string) *string {
	outStr := path_sep_replacer_regex.ReplaceAllString(*str, pathSeparator)
	return &outStr
}

func prepareLine(line *string) (string, string, error) {
	var err error = nil
	outLine := strings.TrimSpace(*line)
	outLine, tag, err := extractTag(&outLine)
	outLine = strings.TrimSpace(outLine)
	outLine = *fixSeparator(&outLine)
	return outLine, tag, err
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func (ignoreList *List) combine(patternList1 *[]pattern, patternList2 *[]pattern) {
	for i := range *patternList2 {
		p2 := &(*patternList2)[i]
		var found bool = false

		for i := range *patternList1 {
			p1 := &(*patternList1)[i]

			if p1.prefix == p2.prefix && p1.suffix == p2.suffix {
				*p1 = *p2
				found = true
				break
			}
		}
		if !found {
			*patternList1 = append(*patternList1, *p2)
		}
	}
}

func (ignoreList *List) hasMatchedPattern(filePath *string, patternList *[]pattern) (bool, int) {
	fixedPath := fixSeparator(filePath)
	for i := range *patternList {
		p := &(*patternList)[i]
		if p.IsEmpty() {
			continue
		}
		if p.isFile {
			if p.prefix == *fixedPath {
				return true, i
			}
		} else {
			if p.HasPrefix() && p.HasSuffix() {
				if strings.HasPrefix(*fixedPath, p.prefix) && strings.HasSuffix(*fixedPath, p.suffix) {
					return true, i
				}
			} else if !p.HasPrefix() && p.HasSuffix() {
				if strings.HasSuffix(*fixedPath, p.suffix) {
					return true, i
				}
			} else if p.HasPrefix() && !p.HasSuffix() {
				if strings.HasPrefix(*fixedPath, p.prefix) {
					return true, i
				}
			}
		}
	}
	return false, -1
}

func (ignoreList *List) processLine(inLine *string) error {
	if len(*inLine) == 0 {
		return nil
	}

	line, tag, err := prepareLine(inLine)
	if err != nil {
		return err
	}

	var actualList *[]pattern
	if strings.HasPrefix(line, not1) || strings.HasPrefix(line, not2) {
		actualList = &ignoreList.includePatternList
	} else {
		actualList = &ignoreList.excludePatternList
	}

	indexLast := strings.LastIndex(line, "*")
	if indexLast != -1 {
		indexFirst := strings.Index(line, "*")
		if indexLast != indexFirst {
			return errors.New(fmt.Sprintf("too many <*> symbols in the pattern <%d>", line))
		}
		list := strings.Split(line, "*")
		*actualList = append(*actualList, pattern{prefix: *removeNot(&list[0]), suffix: list[1], isFile: false, tag: tag})
	} else {
		if strings.HasSuffix(line, pathSeparator) {
			*actualList = append(*actualList, pattern{prefix: *removeNot(&line), isFile: false, tag: tag})
		} else {
			*actualList = append(*actualList, pattern{prefix: *removeNot(&line), isFile: true, tag: tag})
		}
	}
	return nil
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
