package gitlet;

import java.io.Serializable;
import java.io.File;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.TreeMap;

/**
 * @author Felipe Zuluaga
 */

public class Commit implements Serializable {

    /** Given commit's SHA1 ID. */
    private String _commitID;

    /** Message associated to the given commit. */
    private String _message;

    /** Time of given commit. */
    private String _time;

    /** First parent of the given commit. */
    private Commit _parent;

    /** Second parent of the given commit. */
    private Commit _secondParent;

    /** Tree map of files that are being tracked, key would be the
     * fileName, value would be the blob of said file. */
    private TreeMap<String, Blob> _tracking;

    /** Folder where commits are stored. */
    static final File COMMIT_FOLDER = new File(".gitlet/commits/");

    Commit(String msg, long time, Commit parent,
           Commit secondParent, TreeMap<String, Blob> tracking) {
        _message = msg;
        _time = getFormattedTime(time);
        _parent = parent;
        _secondParent = secondParent;
        _tracking = tracking;
        _commitID = ownID();
    }

    public String ownID() {
        byte[] thisCommit = Utils.serialize(this);
        return Utils.sha1(thisCommit);
    }

    public String getFormattedTime(long time) {
        Date date = new Date(time);
        DateFormat formattedDate =
                new SimpleDateFormat("EEE MMM dd HH:mm:ss yyyy");
        String fTime = formattedDate.format(date) + " -0800";
        return fTime;
    }

    public String getID() {
        return _commitID;
    }

    public String getMsg() {
        return _message;
    }

    public String getTime() {
        return _time;
    }

    public Commit getParent() {
        return _parent;
    }

    public Commit getSecondParent() {
        return _secondParent;
    }

    public TreeMap<String, Blob> getTracking() {
        return _tracking;
    }

}
