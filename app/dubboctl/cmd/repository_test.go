// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"testing"

	. "github.com/apache/dubbo-kubernetes/app/dubboctl/internal/testing"
)

// TestRepository_List ensures that the 'list' subcommand shows the client's
// set of repositories by name for builtin repositories, by explicitly
// setting the repositories' path to a new path which includes no others.
func TestRepository_List(t *testing.T) {
	_ = fromTempDirectory(t)

	cmd := NewRepositoryListCmd(NewClient)
	cmd.SetArgs([]string{}) // Do not use test command args

	// Execute the command, capturing the output sent to stdout
	stdout := piped(t)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Assert the output matches expect (whitespace trimmed)
	expect := "default"
	output := stdout()
	if output != expect {
		t.Fatalf("expected:\n'%v'\ngot:\n'%v'\n", expect, output)
	}
}

// TestRepository_Add ensures that the 'add' subcommand accepts its positional
// arguments, respects the repositories' path flag, and the expected name is echoed
// upon subsequent 'list'.
func TestRepository_Add(t *testing.T) {
	url := ServeRepo("repository.git#main", t)
	_ = fromTempDirectory(t)
	t.Log(url)

	var (
		add  = NewRepositoryAddCmd(NewClient)
		list = NewRepositoryListCmd(NewClient)
	)
	// Do not use test command args
	add.SetArgs([]string{})
	list.SetArgs([]string{})

	// add [flags] <old> <new>
	add.SetArgs([]string{
		"newrepo",
		url,
	})

	// Parse flags and args, performing action
	if err := add.Execute(); err != nil {
		t.Fatal(err)
	}

	// List post-add, capturing output from stdout
	if err := list.Execute(); err != nil {
		t.Fatal(err)
	}
}

// TestRepository_Rename ensures that the 'rename' subcommand accepts its
// positional arguments, respects the repositories' path flag, and the name is
// reflected as having been renamed upon subsequent 'list'.
func TestRepository_Rename(t *testing.T) {
	url := ServeRepo("repository.git", t)
	_ = fromTempDirectory(t)

	var (
		add    = NewRepositoryAddCmd(NewClient)
		rename = NewRepositoryRenameCmd(NewClient)
		list   = NewRepositoryListCmd(NewClient)
	)
	// Do not use test command args
	add.SetArgs([]string{})
	rename.SetArgs([]string{})
	list.SetArgs([]string{})

	// add a repo which will be renamed
	add.SetArgs([]string{"newrepo", url})
	if err := add.Execute(); err != nil {
		t.Fatal(err)
	}

	// rename [flags] <old> <new>
	rename.SetArgs([]string{
		"newrepo",
		"renamed",
	})

	// Parse flags and args, performing action
	if err := rename.Execute(); err != nil {
		t.Fatal(err)
	}

	// List post-rename, capturing output from stdout
	if err := list.Execute(); err != nil {
		t.Fatal(err)
	}
}

// TestRepository_Remove ensures that the 'remove' subcommand accepts name as
// its argument, respects the repositories' flag, and the entry is removed upon
// subsequent 'list'.
func TestRepository_Remove(t *testing.T) {
	url := ServeRepo("repository.git", t)
	_ = fromTempDirectory(t)

	var (
		add    = NewRepositoryAddCmd(NewClient)
		remove = NewRepositoryRemoveCmd(NewClient)
		list   = NewRepositoryListCmd(NewClient)
		stdout = piped(t)
	)
	// Do not use test command args
	add.SetArgs([]string{})
	remove.SetArgs([]string{})
	list.SetArgs([]string{})

	// add a repo which will be removed
	add.SetArgs([]string{"newrepo", url})
	if err := add.Execute(); err != nil {
		t.Fatal(err)
	}

	// remove [flags] <name>
	remove.SetArgs([]string{
		"newrepo",
	})

	// Parse flags and args, performing action
	if err := remove.Execute(); err != nil {
		t.Fatal(err)
	}

	// List post-remove, capturing output from stdout
	if err := list.Execute(); err != nil {
		t.Fatal(err)
	}

	// Assert the list output now includes the name from args (whitespace trimmed)
	expect := "Repository added: newrepo\nRepository removed\ndefault"
	output := stdout()
	if output != expect {
		t.Fatalf("expected:\n'%v'\ngot:\n'%v'\n", expect, output)
	}
}
