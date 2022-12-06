package gitlet;

import java.io.File;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Map;
import java.util.TreeMap;
import java.util.List;
import java.util.Arrays;
import java.util.Queue;
import java.util.LinkedList;

/**
 * @author Felipe Zuluaga
 */

public class Commands implements Serializable {

    /** points to the commit that the current branch is on. */
    private Commit _HEAD;

    /** Working Directory. */
    private final File workDir = Main.CWD;

    /** Staging Area Add Directory. */
    private final File addDir = Main.STAGE_ADD;

    /** Staging Area Remove Directory. */
    private final File remDir = Main.STAGE_REMOVE;

    /** Remote directory. */
    private final File remoteDir = Main.REMOTE_FOLDER;

    /** Commit Directory. */
    private final File comDir = Commit.COMMIT_FOLDER;

    /** Blob Directory. */
    private final File blobDir = Blob.BLOB_FOLDER;

    /** Branches Directory. */
    private final File branchDir = Branch.BRANCH_FOLDER;

    /** Head Directory. */
    private final File headDir = Branch.HEAD_FOLDER;

    /** boolean to see if conflict in merge. */
    private boolean conflict = false;

    public void log() {
        Commit commit = getHead();
        while (commit != null) {
            System.out.println("===");
            System.out.println("commit " + commit.getID());
            System.out.println("Date: " + commit.getTime());
            System.out.println(commit.getMsg());
            System.out.println();
            commit = commit.getParent();
        }
    }

    public void gLog() {
        File[] allComms = comDir.listFiles();
        for (File file : allComms) {
            Commit comm = Utils.readObject(file, Commit.class);
            System.out.println("===");
            System.out.println("commit " + comm.getID());
            System.out.println("Date: " + comm.getTime());
            System.out.println(comm.getMsg());
            System.out.println();
        }
    }

    public void status() {
        System.out.println("=== Branches ===");
        List branches = Utils.plainFilenamesIn(branchDir);
        if (branches != null) {
            for (Object branch : branches) {
                String name = branch.toString();
                if (isCurrent(name)) {
                    System.out.println("*" + branch);
                } else {
                    System.out.println(branch);
                }
            }
        }
        System.out.println();
        System.out.println("=== Staged Files ===");
        List staged = Utils.plainFilenamesIn(addDir);
        if (staged != null) {
            for (Object file : staged) {
                System.out.println(file);
            }
        }
        System.out.println();
        System.out.println("=== Removed Files ===");
        List removed = Utils.plainFilenamesIn(remDir);
        if (removed != null) {
            for (Object file : removed) {
                System.out.println(file);
            }
        }
        System.out.println();
        System.out.println("=== Modifications Not Staged For Commit ===");
        for (Map.Entry<String, Blob> file
                : getHead().getTracking().entrySet()) {
            String fileName = file.getKey();
            File fileTest = new File(fileName);
            int result = modifiedCheck(fileTest);
            if (result > 0) {
                if (result == 1) {
                    System.out.println(fileName + " (modified)");
                } else {
                    System.out.println(fileName + " (deleted)");
                }
            }
        }
        System.out.println();
        System.out.println("=== Untracked Files ===");
        List working = Utils.plainFilenamesIn(workDir);
        if (working != null) {
            for (Object file : working) {
                File newFile = new File(file.toString());
                if (untracked(newFile)) {
                    System.out.println(file);
                }
            }
        }
    }

    public void add(String fileName) {
        File checkFile = new File(fileName);
        if (!checkFile.exists()) {
            Main.exitWithError("File does not exist.");
        }

        Blob blob = new Blob(checkFile);
        File idk = new File(blobDir, blob.getName());
        if (idk.exists()) {
            Blob test = Utils.readObject(idk, Blob.class);
            String name = test.getName();
        }

        File stagedAdd = new File(addDir, fileName);
        File stagedRem = new File(remDir, fileName);

        Commit currCom = getHead();

        if (currCom.getTracking().containsKey(fileName)) {
            Blob currComSHA = currCom.getTracking().get(fileName);
            String stagedSHA = blob.getID();
            if (stagedSHA.equals(currComSHA.getID())) {
                stagedAdd.delete();
                stagedRem.delete();
                return;
            }
        }

        File checkRem = new File(remDir, fileName);
        if (checkRem.exists()) {
            checkRem.delete();
        }

        File blobNew = new File(blobDir, blob.getName());
        Utils.writeObject(blobNew, blob);

        byte[] stagedContents = Utils.readContents(checkFile);
        Utils.writeContents(stagedAdd, stagedContents);
    }

    public void rm(String fileName) {
        File checkFile = new File(fileName);
        File stagedRem = new File(remDir, fileName);
        File stagedAdd = new File(addDir, fileName);
        Commit commit = getHead();

        if (!commit.getTracking().containsKey(fileName)
                && !stagedAdd.exists()) {
            Main.exitWithError("No reason to remove the file.");
        }

        if (commit.getTracking().containsKey(fileName)) {
            byte[] stagedContents =
                    commit.getTracking().get(fileName).getBlob();
            Utils.writeContents(stagedRem, stagedContents);
            if (checkFile.exists()) {
                Utils.restrictedDelete(checkFile);
            }
        }

        if (stagedAdd.exists()) {
            stagedAdd.delete();
        }
    }

    public void commit(String msg, Commit secondParent) {
        long time;

        TreeMap<String, Blob> tracking = new TreeMap<String, Blob>();
        File[] filesAdd = addDir.listFiles();
        File[] filesRem = remDir.listFiles();
        File[] filesCom = comDir.listFiles();

        File headFile = new File(headDir, "HEAD");

        if (filesCom.length == 0) {
            Commit initCom = new Commit(msg, 0,
                    null, null, tracking);
            Branch master = new Branch("master", initCom);
            File committed = new File(comDir, initCom.getID());
            File created = new File(branchDir, master.getName());

            Utils.writeObject(committed, initCom);
            Utils.writeObject(created, master);
            Utils.writeObject(headFile, master);

            for (File file : filesAdd) {
                file.delete();
            }
            return;
        }

        if (filesAdd.length + filesRem.length == 0) {
            Main.exitWithError("No changes added to the commit.");
        }

        time = System.currentTimeMillis();
        Commit headCom = getHead();
        tracking = headCom.getTracking();

        if (filesAdd.length != 0) {
            for (File file : filesAdd) {
                tracking.put(file.getName(), new Blob(file));
                file.delete();
            }
        }

        if (filesRem.length != 0) {
            for (File file : filesRem) {
                tracking.remove(file.getName());
                file.delete();
            }
        }

        Commit newCom = new Commit(msg, time, getHead(),
                secondParent, tracking);
        Branch replace = new Branch(getCurrent().getName(), newCom);

        File committed = new File(comDir, newCom.getID());
        File branch = new File(branchDir, getCurrent().getName());

        Utils.writeObject(branch, replace);
        Utils.writeObject(headFile, replace);
        Utils.writeObject(committed, newCom);
    }

    public Commit getHead() {
        File head = new File(headDir, "HEAD");
        Branch headBranch = Utils.readObject(head, Branch.class);
        _HEAD = headBranch.getPoint();
        return _HEAD;
    }

    public Branch getCurrent() {
        File head = new File(headDir, "HEAD");
        Branch headBranch = Utils.readObject(head, Branch.class);
        return headBranch;
    }

    public void setHead(String branchName) {
        File setBranch = new File(branchDir, branchName);
        File headFile = new File(headDir, "HEAD");

        Branch branch = Utils.readObject(setBranch, Branch.class);
        Utils.writeObject(headFile, branch);
    }

    public void reset(String commitID) {
        File branch = new File(branchDir, getCurrent().getName());
        File head = new File(headDir, "HEAD");
        File checkCom = new File(comDir, commitID);
        List working = Utils.plainFilenamesIn(workDir);
        File[] stageAdd = addDir.listFiles();
        File[] stageRem = remDir.listFiles();
        ArrayList<String> allTracked = new ArrayList<String>();

        if (!checkCom.exists()) {
            Main.exitWithError("No commit with that id exists.");
        }

        Commit commit = Utils.readObject(checkCom, Commit.class);

        for (Object file : working) {
            if (!getHead().getTracking().containsKey(file.toString())
                    && commit.getTracking().containsKey(file.toString())) {
                Main.exitWithError("There is an untracked file in the way; "
                        + "delete it, or add and commit it first.");
            }
        }

        TreeMap<String, Blob> tracking = commit.getTracking();

        for (Map.Entry<String, Blob> file
                : tracking.entrySet()) {
            allTracked.add(file.getKey());
        }

        for (String file : allTracked) {
            checkoutCommit(commitID, file);
        }

        for (Object file : working) {
            if (!commit.getTracking().containsKey(file.toString())) {
                File fileDel = new File(file.toString());
                fileDel.delete();
            }
        }

        Branch curr = getCurrent();
        curr.setPoint(commit);
        Utils.writeObject(branch, new Branch(curr.getName(), commit));
        Utils.writeObject(head, new Branch(curr.getName(), commit));

        for (File file : stageAdd) {
            file.delete();
        }
        for (File file : stageRem) {
            file.delete();
        }
    }

    public void find(String msg) {
        List allCommits = Utils.plainFilenamesIn(comDir);
        int found = 0;
        for (Object file : allCommits) {
            File commitFile = new File(comDir, file.toString());
            Commit commit = Utils.readObject(commitFile, Commit.class);
            if (commit.getMsg().equals(msg)) {
                System.out.println(commit.getID());
                found += 1;
            }
        }
        if (found == 0) {
            Main.exitWithError("Found no commit with that message.");
        }
    }

    public void branch(String branchName) {
        File checkBranch = new File(branchDir, branchName);
        if (checkBranch.exists()) {
            Main.exitWithError("A branch with that name already exists.");
        }
        Branch newBranch = new Branch(branchName, getHead());
        File created = new File(branchDir, newBranch.getName());
        Utils.writeObject(created, newBranch);
    }

    public void rmBranch(String branchName) {
        File branch = new File(branchDir, branchName);
        if (!branch.exists()) {
            Main.exitWithError("A branch with that name does not exist.");
        }
        if (isCurrent(branchName)) {
            Main.exitWithError("Cannot remove the current branch.");
        }
        branch.delete();
    }

    public boolean isCurrent(String branchName) {
        String current = getCurrent().getName();
        if (current.equals(branchName)) {
            return true;
        }
        return false;
    }

    public void merge(String branchName) {
        File[] stageAdd = addDir.listFiles();
        File[] stageRem = remDir.listFiles();
        File checkBranch = new File(branchDir, branchName);
        ArrayList<String> untrack = new ArrayList<String>();
        List checkTrack = Utils.plainFilenamesIn(workDir);
        if (stageAdd.length + stageRem.length > 0) {
            Main.exitWithError("You have uncommitted changes.");
        }
        if (!checkBranch.exists()) {
            Main.exitWithError("A branch with that name does not exist.");
        }
        if (getCurrent().getName().equals(branchName)) {
            Main.exitWithError("Cannot merge a branch with itself.");
        }
        Branch givenBranch = Utils.readObject(checkBranch, Branch.class);
        Branch currBranch = getCurrent();
        for (Object file : checkTrack) {
            File check = new File(file.toString());
            if (untracked(check)) {
                Main.exitWithError("There is an untracked file in the way; "
                        + "delete it, or add and commit it first.");
            }
        }
        Commit splitPoint = splitPoint(givenBranch);
        Commit givenCommit = givenBranch.getPoint();
        Commit currentCommit = currBranch.getPoint();
        if (splitPoint.getID().equals(givenCommit.getID())) {
            Main.exitWithError("Given branch is an "
                    + "ancestor of the current branch.");
        }
        if (splitPoint.getID().equals(currentCommit.getID())) {
            checkoutBranch(branchName);
            Main.exitWithError("Current branch fast-forwarded.");
        }
        TreeMap<String, Blob> splitTrack = splitPoint.getTracking();
        TreeMap<String, Blob> givenTrack = givenCommit.getTracking();
        TreeMap<String, Blob> currentTrack = currentCommit.getTracking();
        mergeGiven(givenCommit, currentTrack, splitTrack, givenTrack);
        mergeCurr(currentTrack, splitTrack, givenTrack);
        for (String elem : untrack) {
            currentTrack.remove(elem);
        }
        String mergeMsg = "Merged " + branchName + " into "
                + getCurrent().getName() + ".";
        if (conflict) {
            System.out.println("Encountered a merge conflict.");
        }
        commit(mergeMsg, givenBranch.getPoint());
    }

    public Commit splitPoint(Branch branch) {
        Queue<Commit> ancestorHead = new LinkedList<Commit>();
        Queue<Commit> ancestorGiven = new LinkedList<Commit>();
        Commit headCom = getHead();
        Commit givenCom = branch.getPoint();
        while (givenCom != null) {
            ancestorGiven.add(givenCom);
            if (givenCom.getSecondParent() != null) {
                ancestorGiven.add(givenCom.getSecondParent());
            }
            givenCom = givenCom.getParent();
        }
        while (headCom != null) {
            ancestorHead.add(headCom);
            if (headCom.getSecondParent() != null) {
                ancestorHead.add(headCom.getSecondParent());
            }
            headCom = headCom.getParent();
        }
        for (Commit check : ancestorHead) {
            String currentID = check.getID();
            for (Commit check2 : ancestorGiven) {
                String mergeID = check2.getID();
                if (mergeID.equals(currentID)) {
                    return check;
                }
            }
        }
        return null;
    }

    public void mergeGiven(Commit givenCommit,
                           TreeMap<String, Blob> currentTrack,
                           TreeMap<String, Blob>  splitTrack,
                           TreeMap<String, Blob>  givenTrack) {
        for (Map.Entry<String, Blob> file
                : givenTrack.entrySet()) {
            String fileName = file.getKey();
            byte[] givenContents = file.getValue().getBlob();
            if (splitTrack.containsKey(fileName)
                    && currentTrack.containsKey(fileName)) {
                byte[] splitContents = splitTrack.get(fileName).getBlob();
                byte[] currentContents = currentTrack.get(fileName).getBlob();
                if (Arrays.equals(currentContents, splitContents)
                        && !Arrays.equals(givenContents, splitContents)) {
                    checkoutCommit(givenCommit.getID(), fileName);
                    File staged = new File(addDir, fileName);
                    Utils.writeContents(staged, givenContents);
                }
            }
            if (!currentTrack.containsKey(fileName)
                    && !splitTrack.containsKey(fileName)) {
                checkoutCommit(givenCommit.getID(), fileName);
                File staged = new File(addDir, fileName);
                Utils.writeContents(staged, givenContents);
            }
            if (splitTrack.containsKey(fileName)
                    && !currentTrack.containsKey(fileName)) {
                byte[] splitContents = splitTrack.get(fileName).getBlob();
                if (!Arrays.equals(givenContents, splitContents)) {
                    conflict = true;
                    replaceConflicted(givenContents, null, fileName);
                }
            }
            if (!splitTrack.containsKey(fileName)
                    && currentTrack.containsKey(fileName)) {
                byte[] currContents = currentTrack.get(fileName).getBlob();
                if (!Arrays.equals(givenContents, currContents)) {
                    conflict = true;
                    replaceConflicted(givenContents, currContents, fileName);
                }
            }
        }
    }

    public void mergeCurr(TreeMap<String, Blob> currentTrack,
                          TreeMap<String, Blob>  splitTrack,
                          TreeMap<String, Blob>  givenTrack) {

        for (Map.Entry<String, Blob> file
                : currentTrack.entrySet()) {

            String fileName = file.getKey();
            byte[] currContents = file.getValue().getBlob();
            if (splitTrack.containsKey(fileName)
                    && !givenTrack.containsKey(fileName)) {
                byte[] splitContents = splitTrack.get(fileName).getBlob();
                if (Arrays.equals(currContents, splitContents)) {
                    File removed = new File(fileName);
                    removed.delete();
                }
                if (!Arrays.equals(currContents, splitContents)) {
                    conflict = true;
                    replaceConflicted(null, currContents, fileName);
                }
            }

            if (splitTrack.containsKey(fileName)
                    && givenTrack.containsKey(fileName)) {
                byte[] splitContents = splitTrack.get(fileName).getBlob();
                byte[] givenContents = givenTrack.get(fileName).getBlob();
                if (!Arrays.equals(currContents, splitContents)
                        && !Arrays.equals(givenContents, splitContents)
                        && !Arrays.equals(givenContents, currContents)) {
                    conflict = true;
                    replaceConflicted(givenContents, currContents, fileName);
                }
            }

            if (!splitTrack.containsKey(fileName)
                    && givenTrack.containsKey(fileName)) {
                byte[] givenContents = givenTrack.get(fileName).getBlob();
                if (!Arrays.equals(currContents, givenContents)) {
                    conflict = true;
                    replaceConflicted(givenContents, currContents, fileName);
                }
            }
        }
    }

    public void replaceConflicted(byte[] given,
                              byte[] current, String fileName) {
        String replace = "";
        if (given == null) {
            String currentContents = new String(current);
            replace = "<<<<<<< HEAD" + "\n"
                    + currentContents
                    + "=======" + "\n"
                    + ">>>>>>>" + "\n";
        } else if (current == null) {
            String givenContents = new String(given);
            replace = "<<<<<<< HEAD" + "\n"
                    + "=======" + "\n"
                    + givenContents
                    + ">>>>>>>" + "\n";
        } else {
            String givenContents = new String(given);
            String currentContents = new String(current);
            replace = "<<<<<<< HEAD" + "\n"
                    + currentContents
                    + "=======" + "\n"
                    + givenContents
                    + ">>>>>>>" + "\n";
        }
        File replaceFile = new File(fileName);
        File staged = new File(addDir, fileName);
        Utils.writeContents(replaceFile, replace);
        Utils.writeContents(staged, replace);
    }

    public void checkoutHead(String file) {
        checkoutCommit(getHead().getID(), file);
    }

    public void checkoutCommit(String commitID, String fileName) {
        List allCommits = Utils.plainFilenamesIn(comDir);
        File newFile = new File(fileName);

        for (Object file : allCommits) {
            File commitFile = new File(comDir, file.toString());
            Commit commit = Utils.readObject(commitFile, Commit.class);
            if (commit.getID().startsWith(commitID)) {
                TreeMap<String, Blob> comTrack = commit.getTracking();
                if (!comTrack.containsKey(fileName)) {
                    Main.exitWithError("File does not exist in that commit.");
                }

                byte[] checkout = comTrack.get(fileName).getBlob();
                Utils.writeContents(newFile, checkout);

                return;
            }
        }
        Main.exitWithError("No commit with that id exists.");
    }

    public void checkoutBranch(String branchName) {
        File branches = branchDir;
        File checkBranch = new File(branches, branchName);

        if (!checkBranch.exists()) {
            Main.exitWithError("No such branch exists.");
        } else if (isCurrent(branchName)) {
            Main.exitWithError("No need to checkout the current branch.");
        }

        untrackedErr(checkBranch);
        Branch currBranch = getCurrent();
        Branch givenBranch = Utils.readObject(checkBranch, Branch.class);
        TreeMap<String, Blob> currTrack =
                currBranch.getPoint().getTracking();
        TreeMap<String, Blob> tracking =
                givenBranch.getPoint().getTracking();

        for (Map.Entry<String, Blob> file
                : tracking.entrySet()) {
            File newFile = new File(file.getKey());
            byte[] write = file.getValue().getBlob();
            Utils.writeContents(newFile, write);
        }

        for (Map.Entry<String, Blob> file
                : currTrack.entrySet()) {
            String fileName = file.getKey();
            File newFile = new File(fileName);
            if (!tracking.containsKey(fileName)) {
                newFile.delete();
            }
        }

        setHead(branchName);
    }

    public int modifiedCheck(File file) {
        Commit currCom = getHead();
        String fileName = file.getName();
        File stagedAdd = new File(addDir, fileName);
        File stagedRem = new File(remDir, fileName);

        if (file.exists()) {
            byte[] fileContents = Utils.readContents(file);
            if (currCom.getTracking().containsKey(fileName)) {
                Blob tracking = currCom.getTracking().get(fileName);
                if (!Arrays.equals(fileContents, tracking.getBlob())
                        && (!stagedAdd.exists() && !stagedRem.exists())) {
                    return 1;
                }
            }
            if (stagedAdd.exists()) {
                byte[] addContents = Utils.readContents(stagedAdd);
                if (!Arrays.equals(addContents, fileContents)) {
                    return 1;
                }
            }
        } else {
            if (stagedAdd.exists()) {
                return 2;
            }
            if (!stagedRem.exists()) {
                if (currCom.getSecondParent() != null) {
                    if (currCom.getTracking().containsKey(fileName)
                            && currCom.getSecondParent().getTracking().
                            containsKey(fileName)) {
                        return 2;
                    }
                    return 0;
                }
                if (currCom.getTracking().containsKey(fileName)) {
                    return 2;
                }
            }
        }
        return 0;
    }

    public boolean untracked(File file) {
        String fileName = file.getName();
        File staged = new File(addDir, fileName);
        Commit commit = getHead();
        if (file.exists()) {
            if (!staged.exists()) {
                while (commit != null) {
                    if (commit.getTracking().containsKey(fileName)) {
                        return false;
                    }
                    commit = commit.getParent();
                }
                return true;
            }
        }
        return false;
    }

    public void untrackedErr(File checkBranch) {
        List working = Utils.plainFilenamesIn(workDir);
        Branch currBranch = getCurrent();
        Branch branch = Utils.readObject(checkBranch, Branch.class);
        TreeMap<String, Blob> currTrack = currBranch.getPoint().getTracking();
        TreeMap<String, Blob> tracking = branch.getPoint().getTracking();
        for (Object file : working) {
            if (!currTrack.containsKey(file.toString())
                    && tracking.containsKey(file.toString())) {
                Main.exitWithError("There is an untracked file in the way; "
                        + "delete it, or add and commit it first.");
            }
        }
    }

    public void addRemote(String remoteName, String filePath) {
        String newPath = filePath.replace("/", File.separator);
        File remote = new File(remoteDir, remoteName);
        if (remote.exists()) {
            Main.exitWithError("A remote with that name already exists.");
        }
        File pulledFile = new File(remoteDir, "pulled");
        File countFile = new File(remoteDir, "count");
        boolean pulled = false;
        if (!countFile.exists()) {
            String count = "0";
            Utils.writeContents(countFile, count);
        }
        Utils.writeObject(pulledFile, pulled);
        Utils.writeContents(remote, newPath);
    }

    public void fetch(String remoteName, String branchName) {
        File remote = new File(remoteDir, remoteName);
        if (!remote.exists() || branchName.equals("master")) {
            Main.exitWithError("Remote directory not found.");
        }
        String path = Utils.readContentsAsString(remote);
        String newName = remoteName + "/" + branchName;
        File branch = new File("/branches/", branchName);
        File branchPath = new File(path, branch.toString());
        if (!branchPath.exists()) {
            Main.exitWithError("That remote does not have that branch.");
        }
        Commit pointingTo = Utils.readObject(branchPath,
                Branch.class).getPoint();
        Branch newBranch = new Branch(newName, pointingTo);
        File remoteBranch = new File(branchDir, newName);
        Utils.writeObject(remoteBranch, newBranch);
    }

    public void push(String remoteName, String branchName) {
        File remote = new File(remoteDir, remoteName);
        File pulledFile = new File(remoteDir, "pulled");
        File countFile = new File(remoteDir, "count");
        String count = Utils.readContentsAsString(countFile).toString();
        boolean pulled = Utils.readObject(pulledFile, Boolean.class);
        if (!remote.exists() || branchName.equals("master")
                && count.equals("0")) {
            count = "1";
            Utils.writeContents(countFile, count);
            Main.exitWithError("Remote directory not found.");
        }
        if (pulledFile.exists() && !pulled) {
            Main.exitWithError("Please pull down "
                    + "remote changes before pushing.");
        }
    }

    public void pull(String remoteName, String branchName) {
        fetch(remoteName, branchName);
        String trueBranch = remoteName + "/" + branchName;
        merge(trueBranch);
        File pulledFile = new File(remoteDir, "pulled");
        boolean pulled = true;
        Utils.writeObject(pulledFile, pulled);
    }

    public void rmRemote(String remoteName) {
        File remote = new File(remoteDir, remoteName);
        if (!remote.exists()) {
            Main.exitWithError("A remote with that name does not exist.");
        }
        remote.delete();
    }

}
