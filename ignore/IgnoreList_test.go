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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

const (
	filePath string = "testIgnoreList"
)

func prefixSuffixList() []string {
	return []string{
		"prefix 1/suffix 1",
		"prefix 1/suffix 2",
		"prefix 2/suffix 2",
		"prefix 2/suffix 1",
	}
}

func filePathList() []string {
	return []string{
		"folder1/file1",
		"folder1:file2",
		"folder2/folder1/file1",
		"folder2/folder2\\file1",
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeIgnoreListFile(strList []string) {
	f, err := os.Create(filePath)
	check(err)
	defer f.Close()
	for _, value := range strList {
		f.WriteString(value)
		f.WriteString("\n")
	}
}

func removeIgnoreListFile() {
	err := os.Remove(filePath)
	check(err)
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// *

func TestNoPrefixNoSuffix(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"*", " ", "/", "\\"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.False(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// #*

func TestYesPrefixNoSuffix_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.True(ignoreList.IsIgnored(psList[0]))
	a.True(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.False(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixNoSuffix_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"", "prefix 2*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.False(ignoreList.IsIgnored(psList[0]))
	a.False(ignoreList.IsIgnored(psList[1]))
	a.True(ignoreList.IsIgnored(psList[2]))
	a.True(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixNoSuffix_case5(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*", "prefix 2*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.True(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

func TestYesPrefixNoSuffix_case6(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix*", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.True(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// #*#

func TestYesPrefixYesSuffix_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*suffix 2", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.False(ignoreList.IsIgnored(psList[0]))
	a.True(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.False(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixYesSuffix_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"", "prefix 2*suffix 1"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.False(ignoreList.IsIgnored(psList[0]))
	a.False(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.True(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixYesSuffix_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*suffix 2", "prefix 2*suffix 1"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.False(ignoreList.IsIgnored(psList[0]))
	a.True(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.True(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixYesSuffix_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix*suffix 1", "prefix*suffix 2"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.True(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

func TestYesPrefixYesSuffix_case5(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*suffix 1", "prefix 1*suffix 2"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.True(ignoreList.IsIgnored(psList[0]))
	a.True(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.False(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestYesPrefixYesSuffix_case6(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix 1*suffix 1", "prefix 1*suffix"})
	ignoreList := NewList()
	err := ignoreList.LoadFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.True(ignoreList.IsIgnored(psList[0]))
	a.False(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.False(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// *#

func TestNoPrefixYesSuffix_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"*suffix 1", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.True(ignoreList.IsIgnored(psList[0]))
	a.False(ignoreList.IsIgnored(psList[1]))
	a.False(ignoreList.IsIgnored(psList[2]))
	a.True(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestNoPrefixYesSuffix_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"", "*suffix 2"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	psList := prefixSuffixList()
	a.False(ignoreList.IsIgnored(psList[0]))
	a.True(ignoreList.IsIgnored(psList[1]))
	a.True(ignoreList.IsIgnored(psList[2]))
	a.False(ignoreList.IsIgnored(psList[3]))
	removeIgnoreListFile()
}

func TestNoPrefixYesSuffix_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"*suffix 1", "*suffix 2"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.True(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

func TestNoPrefixYesSuffix_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"*suffix", ""})
	ignoreList := NewList()
	err := ignoreList.LoadFromFile(filePath)
	a.NoError(err)
	for _, value := range prefixSuffixList() {
		a.False(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestPathAddPattern_case1(t *testing.T) {
	a := assert.New(t)
	ignoreList := NewList()
	err := ignoreList.AddPattern("folder2/folder2/file1")
	a.NoError(err)
	err = ignoreList.AddPattern("folder2/folder1/file1")
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3]))  // "folder2/folder2/file1"
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// folder

func TestPathIsFolder_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder", "folder/", "file1", "file1/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range filePathList() {
		a.False(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

func TestPathIsFolder_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1\\", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[1]))  // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"", "folder2/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3]))  // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1/", "folder2/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_case5(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2\\folder1/", ""})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_case6(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2:folder1/", "folder1\\"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[1]))  // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
// file

func TestPathIsFile_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1", "file1", "folder2", "folder1/file1/", "folder2/folder2/file1/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	for _, value := range filePathList() {
		a.False(ignoreList.IsIgnored(value))
	}
	removeIgnoreListFile()
}

func TestPathIsFile_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1:file1", "folder1/file2"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[1]))  // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFile_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2/folder2/file1", "folder2/folder1/file1"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3]))  // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFile_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1/file1", "folder2/folder1:file1"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFile_case5(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1/file1", "folder2\\folder1/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestPathIsFolder_not_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2/*", "not folder2/folder1/file1"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3]))  // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_not_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2/*", "!folder2/folder1/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.True(ignoreList.IsIgnored(fpList[3]))  // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_not_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder2/*", "not folder2/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestPathIsFolder_not_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1/*", "folder2/*", "not folder1/file2", "!folder2/folder2/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))  // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.True(ignoreList.IsIgnored(fpList[2]))  // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

// It tests the situation when folder name and file name have same start
func TestPathFolderFileName_case1(t *testing.T) {
	a := assert.New(t)
	ignoreList := NewList()
	err := ignoreList.AddPattern("folder1/*")
	a.NoError(err)

	str := "folder1/file1"
	a.True(ignoreList.IsIgnored(str))
	str = "folder1-file1"
	a.False(ignoreList.IsIgnored(str))
}

// It tests the situation when folder name and file name have same start
func TestPathFolderFileName_case2(t *testing.T) {
	a := assert.New(t)
	ignoreList := NewList()
	err := ignoreList.AddPattern("folder1-file1")
	a.NoError(err)

	str := "folder1/file2"
	a.False(ignoreList.IsIgnored(str))
	str = "folder1-file1"
	a.True(ignoreList.IsIgnored(str))
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestIncorrectPattern_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1*/file*1", "folde*r2/folder1/*fil*e1"})
	ignoreList, err := NewListFromFile(filePath)
	a.Error(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
	removeIgnoreListFile()
}

func TestIncorrectPattern_case2(t *testing.T) {
	a := assert.New(t)
	ignoreList := NewList()
	err := ignoreList.AddPattern("folder1*/file*1")
	a.Error(err)
	err = ignoreList.AddPattern("folde*r2/folder1/*fil*e1")
	a.Error(err)
	fpList := filePathList()
	a.False(ignoreList.IsIgnored(fpList[0])) // "folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[1])) // "folder1/file2"
	a.False(ignoreList.IsIgnored(fpList[2])) // "folder2/folder1/file1"
	a.False(ignoreList.IsIgnored(fpList[3])) // "folder2/folder2/file1"
}

func TestIncorrectPattern_case3(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"prefix1*/prefix2*", ""})
	_, err := NewListFromFile(filePath)
	a.Error(err)
	removeIgnoreListFile()
}

func TestIncorrectPattern_case4(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"[prefix1*", ""})
	_, err := NewListFromFile(filePath)
	a.Error(err)
	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestCombine_case1(t *testing.T) {
	a := assert.New(t)
	var ignoreList1 List
	var ignoreList2 List
	ps := string(os.PathSeparator)

	ignoreList1.AddPattern("[tag0] A/*A")
	ignoreList1.AddPattern("[tag1] test1/*test2")
	ignoreList1.AddPattern("[tag2] test3/*test3/")

	ignoreList2.AddPattern("[tag0] B/*B")
	ignoreList2.AddPattern("[tag3] test1/*test2/")
	ignoreList2.AddPattern("[tag4] test3/*test3")

	ignoreList1.Combine(&ignoreList2)
	if a.Len(ignoreList1.excludePatternList, 6) {
		a.Equal(pattern{tag: "tag0", prefix: "A" + ps, suffix: "A", isFile: false}, ignoreList1.excludePatternList[0])
		a.Equal(pattern{tag: "tag1", prefix: "test1" + ps, suffix: "test2", isFile: false}, ignoreList1.excludePatternList[1])
		a.Equal(pattern{tag: "tag2", prefix: "test3" + ps, suffix: "test3" + ps, isFile: false}, ignoreList1.excludePatternList[2])
		a.Equal(pattern{tag: "tag0", prefix: "B" + ps, suffix: "B", isFile: false}, ignoreList1.excludePatternList[3])
		a.Equal(pattern{tag: "tag3", prefix: "test1" + ps, suffix: "test2" + ps, isFile: false}, ignoreList1.excludePatternList[4])
		a.Equal(pattern{tag: "tag4", prefix: "test3" + ps, suffix: "test3", isFile: false}, ignoreList1.excludePatternList[5])
	}
}

func TestCombine_case2(t *testing.T) {
	a := assert.New(t)
	var ignoreList1 List
	var ignoreList2 List
	ps := string(os.PathSeparator)

	ignoreList1.AddPattern("[tag0] A/*A")
	ignoreList1.AddPattern("[tag1] test1/*test2/")
	ignoreList1.AddPattern("[tag2] test3/*test3")

	ignoreList2.AddPattern("[tag0] B/*B")
	ignoreList2.AddPattern("[tag3] test1/*test2/")
	ignoreList2.AddPattern("[tag4] test3/*test3")

	ignoreList1.Combine(&ignoreList2)
	if a.Len(ignoreList1.excludePatternList, 4) {
		a.Equal(pattern{tag: "tag0", prefix: "A" + ps, suffix: "A", isFile: false}, ignoreList1.excludePatternList[0])
		a.Equal(pattern{tag: "tag3", prefix: "test1" + ps, suffix: "test2" + ps, isFile: false}, ignoreList1.excludePatternList[1])
		a.Equal(pattern{tag: "tag4", prefix: "test3" + ps, suffix: "test3", isFile: false}, ignoreList1.excludePatternList[2])
		a.Equal(pattern{tag: "tag0", prefix: "B" + ps, suffix: "B", isFile: false}, ignoreList1.excludePatternList[3])
	}
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestClear(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"folder1/", "folder2/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)

	fpList := filePathList()
	a.True(ignoreList.IsIgnored(fpList[0]))
	a.True(ignoreList.IsIgnored(fpList[1]))
	a.True(ignoreList.IsIgnored(fpList[2]))
	a.True(ignoreList.IsIgnored(fpList[3]))

	ignoreList.Clear()
	a.False(ignoreList.IsIgnored(fpList[0]))
	a.False(ignoreList.IsIgnored(fpList[1]))
	a.False(ignoreList.IsIgnored(fpList[2]))
	a.False(ignoreList.IsIgnored(fpList[3]))

	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/

func TestExtractTag_case1(t *testing.T) {
	a := assert.New(t)
	testString := "[path/file"
	_, _, err := extractTag(&testString)
	a.Error(err)
}

func TestExtractTag_case2(t *testing.T) {
	a := assert.New(t)
	testString := "not path/[123]file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal("not path/[123]file", str)
}

func TestExtractTag_case3(t *testing.T) {
	a := assert.New(t)
	testString := "path/file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal("path/file", str)
}

func TestExtractTag_case4(t *testing.T) {
	a := assert.New(t)
	testString := ""
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal("", str)
}

//--------------------------------------------------------------------------//

func TestExtractTag_case5(t *testing.T) {
	a := assert.New(t)
	testString := "[]"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal("", str)
}

func TestExtractTag_case6(t *testing.T) {
	a := assert.New(t)
	testString := "[]path/file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal("path/file", str)
}

func TestExtractTag_case7(t *testing.T) {
	a := assert.New(t)
	testString := "[] path/file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("", tag)
	a.Equal(" path/file", str)
}

//--------------------------------------------------------------------------//

func TestExtractTag_case8(t *testing.T) {
	a := assert.New(t)
	testString := "[test]"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("test", tag)
	a.Equal("", str)
}

func TestExtractTag_case9(t *testing.T) {
	a := assert.New(t)
	testString := "[test]path/file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("test", tag)
	a.Equal("path/file", str)
}

func TestExtractTag_case10(t *testing.T) {
	a := assert.New(t)
	testString := "[test] path/file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("test", tag)
	a.Equal(" path/file", str)
}

func TestExtractTag_case11(t *testing.T) {
	a := assert.New(t)
	testString := "[test]path/[123]file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("test", tag)
	a.Equal("path/[123]file", str)
}

func TestExtractTag_case12(t *testing.T) {
	a := assert.New(t)
	testString := "[test] path/[123]file"
	str, tag, err := extractTag(&testString)
	a.NoError(err)
	a.Equal("test", tag)
	a.Equal(" path/[123]file", str)
}

//--------------------------------------------------------------------------//

func TestTagWhenExclude_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"[test1]folder1/", "[test2] folder2/"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()

	res, tag := ignoreList.IsIgnoredEx(fpList[0]) // "folder1/file1"
	a.True(res)
	a.Equal("test1", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[1]) // "folder1/file2"
	a.True(res)
	a.Equal("test1", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[2]) // "folder2/folder1/file1"
	a.True(res)
	a.Equal("test2", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[3]) // "folder2/folder2/file1"
	a.True(res)
	a.Equal("test2", tag)

	removeIgnoreListFile()
}

//--------------------------------------------------------------------------//

func TestTagWhenInclude_case1(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"[test1] folder2/*", "[test2] not folder2/folder1/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)
	fpList := filePathList()

	res, tag := ignoreList.IsIgnoredEx(fpList[0]) // "folder1/file1"
	a.False(res)
	a.Equal("", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[1]) // "folder1/file2"
	a.False(res)
	a.Equal("", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[2]) // "folder2/folder1/file1"
	a.False(res)
	a.Equal("test2", tag)

	res, tag = ignoreList.IsIgnoredEx(fpList[3]) // "folder2/folder2/file1"
	a.True(res)
	a.Equal("test1", tag)

	removeIgnoreListFile()
}

// Gets the tag even the file is in include pattern list only
func TestTagWhenInclude_case2(t *testing.T) {
	a := assert.New(t)
	writeIgnoreListFile([]string{"[test2] not folder2/folder1/*"})
	ignoreList, err := NewListFromFile(filePath)
	a.NoError(err)

	file := "folder2/folder1/file1"
	res, tag := ignoreList.IsIgnoredEx(file)
	a.False(res)
	a.Equal("test2", tag)

	removeIgnoreListFile()
}

/*********************************************************************************************************/
///////////////////////////////////////////////////////////////////////////////////////////////////////////
/*********************************************************************************************************/
