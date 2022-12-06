package gitlet;

import java.io.File;

/** Driver class for Gitlet, the tiny stupid version-control system.
 *  @author Felipe Zuluaga
 */
public class Main {
    /** Current Working Directory. */
    static final File CWD = new File(".");

    /** Main metadata folder. */
    static final File GIT_FOLDER =
            new File(".gitlet");

    /** remote directory. */
    static final File REMOTE_FOLDER =
            new File("../remote");

    /** Staging Area directory. */
    static final File STAGING_AREA =
            new File(".gitlet/stagingArea");

    /** Staging Area Add folder. */
    static final File STAGE_ADD =
            new File(".gitlet/stagingArea/stageAdd");

    /** Staging Area Remove folder. */
    static final File STAGE_REMOVE =
            new File(".gitlet/stagingArea/stageRemove");


    /** Usage: java gitlet.Main ARGS, where ARGS contains.
     *  <COMMAND> <OPERAND> ....
     *
     */
    public static void main(String[] args) {
        if (args.length == 0) {
            System.out.println("Please enter a command.");
            return;
        }
        if (!args[0].equals("init") && !initialized()) {
            exitWithError("Not in an initialized Gitlet directory.");
        }
        Commands c = new Commands();
        switch (args[0]) {
        case "init" :
            if (validateNumArgs(args, 1)) {
                init();
            }
            return;
        case "log" :
            if (validateNumArgs(args, 1)) {
                c.log();
            }
            return;
        case "global-log" :
            if (validateNumArgs(args, 1)) {
                c.gLog();
            }
            break;
        case "status" :
            if (validateNumArgs(args, 1)) {
                c.status();
            }
            break;
        case "add" :
            if (validateNumArgs(args, 2)) {
                String fileName = args[1];
                c.add(fileName);
            }
            return;
        case "rm" :
            if (validateNumArgs(args, 2)) {
                String fileName = args[1];
                c.rm(fileName);
            }
            return;
        case "commit" :
            if (validateNumArgs(args, 2)) {
                if (args[1].isEmpty()) {
                    exitWithError("Please enter a commit message.");
                }
                String msg = args[1];
                c.commit(msg, null);
            }
            return;
        default:
            main2(args, c);
        }
    }

    public static void main2(String[] args, Commands c) {
        switch (args[0]) {
        case "find" :
            if (validateNumArgs(args, 2)) {
                String commitID = args[1];
                c.find(commitID);
            }
            return;
        case "reset" :
            if (validateNumArgs(args, 2)) {
                String msg = args[1];
                c.reset(msg);
            }
            return;
        case "checkout" :
            switch (args[1]) {
            case "--" :
                if (validateNumArgs(args, 3)) {
                    String file = args[2];
                    c.checkoutHead(file);
                }
                return;
            default:
                if (args.length == 2) {
                    String branch = args[1];
                    c.checkoutBranch(branch);
                } else if (validateNumArgs(args, 4)) {
                    if (args[2].equals("--")) {
                        String commitID = args[1];
                        String file = args[3];
                        c.checkoutCommit(commitID, file);
                        return;
                    }
                    exitWithError("Incorrect operands.");
                }
                return;
            }
        case "branch" :
            if (validateNumArgs(args, 2)) {
                String branch = args[1];
                c.branch(branch);
            }
            return;
        case "rm-branch" :
            if (validateNumArgs(args, 2)) {
                String branch = args[1];
                c.rmBranch(branch);
            }
            return;
        case "merge" :
            if (validateNumArgs(args, 2)) {
                String branch = args[1];
                c.merge(branch);
            }
            return;
        default:
            mainRemote(args, c);
        }
    }

    public static void mainRemote(String[] args, Commands c) {
        switch (args[0]) {
        case "add-remote":
            if (validateNumArgs(args, 3)) {
                c.addRemote(args[1], args[2]);
            }
            return;
        case "rm-remote":
            if (validateNumArgs(args, 2)) {
                c.rmRemote(args[1]);
            }
            return;
        case "fetch":
            if (validateNumArgs(args, 3)) {
                c.fetch(args[1], args[2]);
            }
            return;
        case "push":
            if (validateNumArgs(args, 3)) {
                c.push(args[1], args[2]);
            }
            return;
        case "pull":
            if (validateNumArgs(args, 3)) {
                c.pull(args[1], args[2]);
            }
            return;
        default:
            exitWithError("No command with that name exists.");
        }
    }


    public static void init() {
        Commands c = new Commands();
        if (!GIT_FOLDER.exists()) {
            GIT_FOLDER.mkdir();
            if (!REMOTE_FOLDER.exists()) {
                REMOTE_FOLDER.mkdir();
            }
            STAGING_AREA.mkdir();
            STAGE_ADD.mkdir();
            STAGE_REMOVE.mkdir();
            Branch.BRANCH_FOLDER.mkdir();
            Branch.HEAD_FOLDER.mkdir();
            Commit.COMMIT_FOLDER.mkdir();
            Blob.BLOB_FOLDER.mkdir();
            c.commit("initial commit", null);
            return;
        }
        exitWithError("A Gitlet version-control system already exists "
                + "in the current directory.");
    }

    public static void exitWithError(String message) {
        if (message != null && !message.equals("")) {
            System.out.println(message);
        }
        System.exit(0);
    }

    public static boolean validateNumArgs(String[] args, int n) {
        if (args.length != n) {
            exitWithError("Incorrect operands.");
        }
        return true;
    }

    public static boolean initialized() {
        if (!GIT_FOLDER.exists()) {
            return false;
        }
        return true;
    }

}
