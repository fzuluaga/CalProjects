package enigma;

/** An alphabet of encodable characters.  Provides a mapping from characters
 *  to and from indices into the alphabet.
 *  @author Felipe Zuluaga
 */
class Alphabet {
    /** A new alphabet containing CHARS. The K-th character has index
     *  K (numbering from 0). No character may be duplicated. */
    Alphabet(String chars) {
        _alpha = new char[chars.length()];
        _length = _alpha.length;
        for (int i = 0; i < chars.length(); i++) {
            _alpha[i] = chars.charAt(i);
        }
    }

    /** A default alphabet of all upper-case characters. */
    Alphabet() {
        this("ABCDEFGHIJKLMNOPQRSTUVWXYZ");
    }

    /** Returns the size of the alphabet. */
    int size() {
        return _length;
    }

    /** Returns true if CH is in this alphabet. */
    boolean contains(char ch) {
        for (int i = 0; i < size(); i++) {
            if (_alpha[i] == ch) {
                return true;
            }
        }
        return false;
    }

    /** Returns character number INDEX in the alphabet, where
     *  0 <= INDEX < size(). */
    char toChar(int index) {
        return _alpha[index];
    }

    /** Returns the index of character CH which must be in
     *  the alphabet. This is the inverse of toChar(). */
    int toInt(char ch) {
        if (contains(ch)) {
            for (int i = 0; i < size(); i++) {
                if (_alpha[i] == ch) {
                    return i;
                }
            }
        }
        return 0;
    }

    /** Length of the alphabet. */
    private int _length;

    /** A CharList containing all the chars in the alphabet. */
    private char[] _alpha;

}
