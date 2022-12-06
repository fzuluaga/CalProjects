package gitlet;

import java.io.File;
import java.io.Serializable;

/**
 * @author Felipe Zuluaga
 */

public class Blob implements Serializable {

    /** Name of given blob. */
    private String _name;

    /** Contents of given blob. */
    private byte[] _blob;

    /** SHA1 ID of given blob. */
    private String _id;

    /** Folder where blobs are stored. */
    static final File BLOB_FOLDER = new File(".gitlet/blobs");

    Blob(File file) {
        _blob = Utils.readContents(file);
        _name = file.getName();
        _id = Utils.sha1(this._blob);
    }

    public String getID() {
        return _id;
    }

    public byte[] getBlob() {
        return _blob;
    }

    public String getName() {
        return _name;
    }

}
