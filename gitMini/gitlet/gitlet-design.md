# Gitlet Design Document
author: Felipe Zuluaga


## Classes and Data Structures

Blob : The contents of the file being manipulated (Object).

Commit History : A tree with various nodes that point to commits and 
various branching paths that it has gone through (Tree).

Commits: A list of commits that point to a blob and contain its own metadata,
(Linked List).

### Main

This class is the main class of Gitlet, will be used for functions 
that are a commonality between all other functions, will also instantiate
and hold all necessary instance variables.

### Blob

This class will instantiate blobs and set and edit the information within them.

Instance Variables:
* int _id
* File _file

### FileStaging

Stages and unstages files to be commited as neccessary in order for them
to be commited.

Instance Variables:
* boolean _added 

### Commit

Performs a commit by setting combination of log messages, timestamp, 
a reference to a tree, and references to parent commits.

Instance Variables:
* int _id
* int _timestamp
* String _message
* String _parent
* String _secondParent ( if applicable )
* ArrayList _tracking

### Branching

Takes care of merging and the addition/removal of branches in the commit history.

Instance Variables:
* Commit _parent
* Commit _secondParent
* Commit _head

## Algorithms

### Gitlet Class
1. init(): Initializes a gitlet control system in the CWD. starts with an empty
commit with the message "initial commit". Also is initialized with a single branch
called master. UID for all starting commits must be the same as they will have identical
content. Instance variables will have to be set and commits and branches will be made.
2. main(): check the arguments of the command being inputted and then act accordingly.
will have a switch case format checking the argument of the command and then running
the proper function.
3. status(): Displays the current status of the gitlet, shows what branches exist
and which one we are currently on, as well as displaying what files are not being
tracked, have been staged, removed, and modified. 

### File Staging Class
1. add(): Add a file to the folder stagingArea in .gitlet this will then allow for it
to be commited, will only work if the CWV is not the same as the file being staged. 
2. rm(): Will remove the file ine the stagingArea, only remove it if it is being tracked
in the current commit.

### Branching Class
1. branch(): creates a new branch with a given name and points it at HEAD node.
2. rmBranch(): removes the pointer associated to the branch with the given name.
3. merge()
4. checkoutCommit(): takes the version of the file in the commit of the given id, 
puts it in the CWD, overwrites the current present file and does not stage the file.
5. checkoutBranch(): Does the same as checkoutCommit but for a branch, takes all
files in the given branch and puts it in the CWD overwriting the versions of the
files that are present in it already. Changes the HEAD to the current branch all files
that are present in current branch that aren't in the checked out branch are deleted
and staging area is cleared.

work done for checkout can be simplified and reused in both functions.

### Commit Class
1. commit(): Check if files are present in the staging area and then commit them.
2. log(): displays information about each commit in backward order. shows commit
id, timestamp, and commit message.
3. globalLog(): same as log, but shows all commits ever made.
4. find(): Finds and prints out the commit id of all (if multiple) commits with the
given commit message. Will have parameters commitMessage.
5. reset(): Checks out all files that are being tracked by a given commit, and staging
Area is cleared.


## Persistence

In order to ensure that the files that are commited and added are persistent
the blobs, and the state of the commit nodes will have to be saved after each
call to the git machine. To do this,

1. Write all the Blob objects to disk. Serialize all the Blob objects that are 
   being used and write them into files on the disk (using the blob SHA-1 code)
   This can, and will, be done by using the writeObject method from the Utils class. 
   The Blob class will have to be Serializable.


2. Write all Commit objects and their metadate to disk. Serialize all the commits
   into bytes that can then be written onto a specific file on the disk, then do the
   same thing that is done for Blob objects to the metadata of the commit. These can
   both be done with the writeObject method from the Utils class.


3. Write all Commit Branch History to disk. Serialize the commit tree into bytes that
   that can then be written onto a specifically named file using the commit hash on 
   the disk. This can also be done using the writeObject Method from the Utils class.
   Will be useful for globalLog().

### Gitlet Directory:

* .gitlet/ :
  * stagingArea/
    * stageAdd
    * stageRemove
  * commits/
    * .commitHistory ? 