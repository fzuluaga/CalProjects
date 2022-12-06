package enigma;

import static enigma.EnigmaException.*;

/** Class that represents a rotating rotor in the enigma machine.
 *  @author Felipe Zuluaga
 */
class MovingRotor extends Rotor {

    /** A rotor named NAME whose permutation in its default setting is
     *  PERM, and whose notches are at the positions indicated in NOTCHES.
     *  The Rotor is initally in its 0 setting (first character of its
     *  alphabet).
     */
    MovingRotor(String name, Permutation perm, String notches) {
        super(name, perm);
        _notches = notches;
    }

    @Override
    boolean rotates() {
        return true;
    }

    @Override
    void advance() {
        int nextSetting = permutation().wrap(setting() + 1);
        super.set(nextSetting);
    }

    @Override
    boolean atNotch() {
        for (int i = 0; i < notches().length(); i++) {
            if (notches().charAt(i) == alphabet().toChar(setting())) {
                return true;
            }
        }
        return false;
    }

    @Override
    String notches() {
        String notches = "";
        for (int i = 0; i < _notches.length(); i++) {
            int diff = alphabet().toInt(_notches.charAt(i)) - ring();
            int indexNotch = permutation().wrap(diff);
            notches += alphabet().toChar(indexNotch);
        }
        return notches;
    }

    /** Notches of the rotor. */
    private String _notches;

}
