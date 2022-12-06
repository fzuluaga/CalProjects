package gitlet;

import java.io.File;
import java.io.Serializable;

/**
 * @author Felipe Zuluaga
 */

public class Branch implements Serializable {

    /** Branches directory. */
    static final File BRANCH_FOLDER = new File(".gitlet/branches");

    /** Head branch directory. */
    static final File HEAD_FOLDER = new File(".gitlet/head");

    /** Commit that branch is pointing to. */
    private Commit _pointing;

    /** Name of given branch. */
    private String _branchName;

    Branch(String name, Commit pointTo) {
        _pointing = pointTo;
        _branchName = name;
    }

    public Commit getPoint() {
        return _pointing;
    }

    public void setPoint(Commit commit) {
        _pointing = commit;
    }

    public String getName() {
        return _branchName;
    }

}
