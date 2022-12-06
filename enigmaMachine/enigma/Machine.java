package enigma;

import java.util.Collection;

import static enigma.EnigmaException.*;

/** Class that represents a complete enigma machine.
 *  @author Felipe Zuluaga
 */
class Machine {

    /** A new Enigma machine with alphabet ALPHA, 1 < NUMROTORS rotor slots,
     *  and 0 <= PAWLS < NUMROTORS pawls.  ALLROTORS contains all the
     *  available rotors. */
    Machine(Alphabet alpha, int numRotors, int pawls,
            Collection<Rotor> allRotors) {
        _alphabet = alpha;
        _pawls = pawls;
        _numRotors = numRotors;
        _allRotors = allRotors.toArray();
        _currentRotors = new Rotor[_numRotors];

    }

    /** Return the number of rotor slots I have. */
    int numRotors() {
        return _numRotors;
    }

    /** Return the number pawls (and thus rotating rotors) I have. */
    int numPawls() {
        return _pawls;
    }

    /** Return Rotor #K, where Rotor #0 is the reflector, and Rotor
     *  #(numRotors()-1) is the fast Rotor.  Modifying this Rotor has
     *  undefined results. */
    Rotor getRotor(int k) {
        return _currentRotors[k];
    }

    Alphabet alphabet() {
        return _alphabet;
    }

    /** Set my rotor slots to the rotors named ROTORS from my set of
     *  available rotors (ROTORS[0] names the reflector).
     *  Initially, all rotors are set at their 0 setting. */
    void insertRotors(String[] rotors) {
        for (int i = 0; i < rotors.length; i++) {
            for (int j = 0; j < _allRotors.length; j++) {
                if (rotors[i].equals(((Rotor) _allRotors[j]).name())) {
                    _currentRotors[i] = (Rotor) _allRotors[j];
                    break;
                }
            }
        }
    }


    /** Set my rotors according to SETTING, which must be a string of
     *  numRotors()-1 characters in my alphabet. The first letter refers
     *  to the leftmost rotor setting (not counting the reflector).
     *  @param ringstellung string of the ringstellung.
     *  @param setting string of the setting.
     *  */
    void setRotors(String setting,  String ringstellung) {
        if (!ringstellung.equals("AAAA")) {
            for (int i = 1; i < _currentRotors.length; i++) {
                int set1 = _alphabet.toInt(setting.charAt(i - 1));
                int ring1 = _alphabet.toInt(ringstellung.charAt(i - 1));
                int diff = _currentRotors[0].permutation().wrap(set1 - ring1);
                _currentRotors[i].setRing(ring1);
                _currentRotors[i].set(diff);
            }

        } else {
            for (int i = 1; i < _currentRotors.length; i++) {
                _currentRotors[i].set(setting.charAt(i - 1));
                _currentRotors[i].setRing(0);
            }
        }
    }

    /** Return the current plugboard's permutation. */
    Permutation plugboard() {
        return _plugboard;
    }

    /** Set the plugboard to PLUGBOARD. */
    void setPlugboard(Permutation plugboard) {
        _plugboard = plugboard;
    }

    /** Returns the result of converting the input character C (as an
     *  index in the range 0..alphabet size - 1), after first advancing
     *  the machine. */
    int convert(int c) {
        advanceRotors();
        if (Main.verbose()) {
            System.err.printf("[");
            for (int r = 1; r < numRotors(); r += 1) {
                System.err.printf("%c",
                        alphabet().toChar(getRotor(r).setting()));
            }
            System.err.printf("] %c -> ", alphabet().toChar(c));
        }
        c = plugboard().permute(c);
        if (Main.verbose()) {
            System.err.printf("%c -> ", alphabet().toChar(c));
        }
        c = applyRotors(c);
        c = plugboard().permute(c);
        if (Main.verbose()) {
            System.err.printf("%c%n", alphabet().toChar(c));
        }
        return c;
    }

    /** Advance all rotors to their next position. */
    private void advanceRotors() {
        boolean[] willAdvance = new boolean[_currentRotors.length];
        for (int i = 0; i < _currentRotors.length; i++) {
            if (!_currentRotors[i].rotates()) {
                willAdvance[i] = false;
            } else if (i == _currentRotors.length - 1) {
                willAdvance[i] = true;
            } else if (_currentRotors[i + 1].atNotch()) {
                willAdvance[i] = true;
                willAdvance[i + 1] = true;
            }
        }
        for (int i = 0; i < willAdvance.length; i++) {
            if (willAdvance[i]) {
                _currentRotors[i].advance();
            }
        }
    }

    /** Return the result of applying the rotors to the character C (as an
     *  index in the range 0..alphabet size - 1). */
    private int applyRotors(int c) {
        boolean reverse = false;
        for (int i = _currentRotors.length - 1; i >= 0; i--) {
            if (!reverse) {
                c = _currentRotors[i].convertForward(c);
                if (_currentRotors[i].reflecting()) {
                    reverse = true;
                }
            }
        }
        for (int i = 1; i < _currentRotors.length; i++) {
            c = _currentRotors[i].convertBackward(c);
        }
        return c;
    }

    /** Returns the encoding/decoding of MSG, updating the state of
     *  the rotors accordingly. */
    String convert(String msg) {
        String result = "";
        for (int i = 0; i < msg.length(); i++) {
            int charInd = _alphabet.toInt(msg.charAt(i));
            char convertedChar = _alphabet.toChar(convert(charInd));
            result += convertedChar;
        }
        return result;
    }

    /** Common alphabet of my rotors. */
    private final Alphabet _alphabet;

    /** Number of rotors in machine. */
    private int _numRotors;

    /** Number of pawls in machine. */
    private int _pawls;

    /** List of all available rotors. */
    private Object[] _allRotors;

    /** The rotors currently being used. */
    private Rotor[] _currentRotors;

    /** Permutation of plugboard. */
    private Permutation _plugboard;

}
