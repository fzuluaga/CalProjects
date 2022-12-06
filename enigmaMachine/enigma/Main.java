package enigma;

import java.io.File;
import java.io.IOException;
import java.io.PrintStream;

import java.util.Scanner;
import java.util.ArrayList;
import java.util.List;
import java.util.NoSuchElementException;

import ucb.util.CommandArgs;

import static enigma.EnigmaException.*;

/** Enigma simulator.
 *  @author Felipe Zuluaga
 */
public final class Main {

    /** Process a sequence of encryptions and decryptions, as
     *  specified by ARGS, where 1 <= ARGS.length <= 3.
     *  ARGS[0] is the name of a configuration file.
     *  ARGS[1] is optional; when present, it names an input file
     *  containing messages.  Otherwise, input comes from the standard
     *  input.  ARGS[2] is optional; when present, it names an output
     *  file for processed messages.  Otherwise, output goes to the
     *  standard output. Exits normally if there are no errors in the input;
     *  otherwise with code 1. */
    public static void main(String... args) {
        try {
            CommandArgs options =
                new CommandArgs("--verbose --=(.*){1,3}", args);
            if (!options.ok()) {
                throw error("Usage: java enigma.Main [--verbose] "
                            + "[INPUT [OUTPUT]]");
            }

            _verbose = options.contains("--verbose");
            new Main(options.get("--")).process();
            return;
        } catch (EnigmaException excp) {
            System.err.printf("Error: %s%n", excp.getMessage());
        }
        System.exit(1);
    }

    /** Open the necessary files for non-option arguments ARGS (see comment
      *  on main). */
    Main(List<String> args) {
        _config = getInput(args.get(0));

        if (args.size() > 1) {
            _input = getInput(args.get(1));
        } else {
            _input = new Scanner(System.in);
        }

        if (args.size() > 2) {
            _output = getOutput(args.get(2));
        } else {
            _output = System.out;
        }
    }

    /** Return a Scanner reading from the file named NAME. */
    private Scanner getInput(String name) {
        try {
            return new Scanner(new File(name));
        } catch (IOException excp) {
            throw error("could not open %s", name);
        }
    }

    /** Return a PrintStream writing to the file named NAME. */
    private PrintStream getOutput(String name) {
        try {
            return new PrintStream(new File(name));
        } catch (IOException excp) {
            throw error("could not open %s", name);
        }
    }

    /** Configure an Enigma machine from the contents of configuration
     *  file _config and apply it to the messages in _input, sending the
     *  results to _output. */
    private void process() {
        Machine machine = readConfig();
        String setting = _input.nextLine();
        setUp(machine, setting);
        while (_input.hasNextLine()) {
            String nextLine = _input.nextLine();
            if (nextLine.contains("*")) {
                setUp(machine, nextLine);
            } else {
                String secretMessage = nextLine.replaceAll(" ", "");
                String convertedMessage = machine.convert(secretMessage);
                printMessageLine(convertedMessage);
            }
        }
    }

    /** Return an Enigma machine configured from the contents of configuration
     *  file _config. */
    private Machine readConfig() {
        try {
            String alpha = _config.next();
            _alphabet = new Alphabet(alpha);
            _numRotors = _config.nextInt();
            _pawls = _config.nextInt();
            while (_config.hasNext()) {
                _allRotors.add(readRotor());
            }
            return new Machine(_alphabet, _numRotors, _pawls, _allRotors);
        } catch (NoSuchElementException excp) {
            throw error("configuration file truncated");
        }
    }

    /** Return a rotor, reading its description from _config. */
    private Rotor readRotor() {
        try {
            String name = _config.next();
            String type = _config.next();
            String permCycles = "";
            String notches = "";

            while (_config.hasNext("\\(.+\\)")) {
                permCycles += _config.next();
            }

            Permutation perm = new Permutation(permCycles, _alphabet);

            if (type.charAt(0) == 'M') {
                for (int i = 1; i < type.length(); i++) {
                    notches += type.charAt(i);
                }
                return new MovingRotor(name, perm, notches);
            } else if (type.charAt(0) == 'R') {
                return new Reflector(name, perm);
            } else if (type.charAt(0) == 'N') {
                return new FixedRotor(name, perm);
            } else {
                throw new EnigmaException("No Rotor Present.");
            }
        } catch (NoSuchElementException excp) {
            throw error("bad rotor description");
        }
    }

    /** Set M according to the specification given on SETTINGS,
     *  which must have the format specified in the assignment. */
    private void setUp(Machine M, String settings) {
        Scanner settingLine = new Scanner(settings);
        Scanner arraySet = new Scanner(settings);
        ArrayList<String> usedRotors1 = new ArrayList<>();
        ArrayList<String> pawlsList = new ArrayList<>();
        ArrayList<String> nonReflectors = new ArrayList<>();
        String setting = ""; String plugboard = ""; String ring = "AAAA";
        if (settingLine.next().equals("*")) {
            while (arraySet.hasNext()) {
                String check = arraySet.next();
                if (checkRotor(check)) {
                    usedRotors1.add(check); settingLine.next();
                }
            }
            String[] usedRotors = usedRotors1.toArray(new String[0]);
            checkReflectorPlace(usedRotors1, usedRotors);
            checkRotorRepeat(usedRotors, 0); checkRotorName(usedRotors);
            for (int i = 0; i < usedRotors.length; i++) {
                if (!checkReflects(usedRotors1.get(i))
                        && !checkFixed(usedRotors1.get(i))) {
                    pawlsList.add(usedRotors1.get(i));
                }
                if (!checkReflects(usedRotors1.get(i))) {
                    nonReflectors.add(usedRotors1.get(i));
                }
            }
            if (settingLine.hasNext()) {
                String afterPawls = settingLine.next();
                if (pawlsList.size() > _pawls || checkRotor(afterPawls)) {
                    throw new EnigmaException("Wrong number of arguments.");
                }
                if (!afterPawls.matches("\\(.+\\)")) {
                    setting = afterPawls;
                    for (int i = 0; i < setting.length(); i++) {
                        if (!_alphabet.contains(setting.charAt(i))) {
                            throw new EnigmaException("Bad character.");
                        }
                    }
                    if (setting.length() < nonReflectors.size()) {
                        throw new EnigmaException("Wheel settings too short.");
                    } else if (setting.length() > nonReflectors.size()) {
                        throw new EnigmaException("Wheel settings too long.");
                    }
                }
            }
            if (settingLine.hasNext() && !settingLine.hasNext("\\(.+\\)")) {
                String afterSetting = settingLine.next();
                ring = setUpRing(ring, afterSetting, pawlsList, nonReflectors);
            }
            while (settingLine.hasNext("\\(.+\\)")) {
                plugboard += settingLine.next();
            }
            Permutation plugPerm = new Permutation(plugboard, _alphabet);
            M.setPlugboard(plugPerm); M.insertRotors(usedRotors);
            M.setRotors(setting, ring);
        } else {
            throw new EnigmaException("Missing settings.");
        }
    }

    private String setUpRing(String ring, String afterSetting,
                           ArrayList pList, ArrayList nRList) {
        if (pList.size() > _pawls || checkRotor(afterSetting)) {
            throw new EnigmaException("Wrong number of arguments.");
        }
        if (!afterSetting.matches("\\(.+\\)")) {
            ring = afterSetting;
            for (int i = 0; i < ring.length(); i++) {
                if (!_alphabet.contains(ring.charAt(i))) {
                    throw new EnigmaException("Bad character.");
                }
            }
            if (ring.length() < nRList.size()) {
                throw new EnigmaException("Wheel settings too short.");
            } else if (ring.length() > nRList.size()) {
                throw new EnigmaException("Wheel settings too long.");
            }
        }
        return ring;
    }

    private boolean checkRotor(String rotor) {
        for (int j = 0; j < _allRotors.size(); j++) {
            if (_allRotors.get(j).name().equals(rotor)) {
                return true;
            }
        }
        return false;
    }

    private void checkRotorRepeat(String[] usedRotors, int appearCount) {
        for (int i = 0; i < usedRotors.length; i++) {
            for (int j = 0; j < usedRotors.length; j++) {
                if (usedRotors[i].equals(usedRotors[j])) {
                    appearCount += 1;
                }
            }
            if (appearCount > 1) {
                throw new EnigmaException("Repeated rotor.");
            }
            appearCount = 0;
        }
    }

    private void checkRotorName(String[] usedRotors) {
        ArrayList<String> rotorNames = new ArrayList<>(_allRotors.size());
        for (int j = 0; j < _allRotors.size(); j++) {
            String name = _allRotors.get(j).name();
            rotorNames.add(name);
        }
        for (int i = 0; i < usedRotors.length; i++) {
            if (!rotorNames.contains(usedRotors[i])) {
                throw new EnigmaException("Bad rotor name.");
            }
        }
    }

    private void checkReflectorPlace(ArrayList<String> usedRotors1,
                                     String[] usedRotors) {
        if (!checkReflects(usedRotors[0])) {
            throw new EnigmaException("Reflector in wrong place.");
        }
    }

    private boolean checkReflects(String rotor) {
        for (int j = 0; j < _allRotors.size(); j++) {
            if (_allRotors.get(j).name().equals(rotor)) {
                if (_allRotors.get(j).reflecting()) {
                    return true;
                }
            }
        }
        return false;
    }

    private boolean checkFixed(String rotor) {
        for (int j = 0; j < _allRotors.size(); j++) {
            if (_allRotors.get(j).name().equals(rotor)) {
                if (!_allRotors.get(j).rotates()) {
                    return true;
                }
            }
        }
        return false;
    }

    /** Return true iff verbose option specified. */
    static boolean verbose() {
        return _verbose;
    }

    /** Print MSG in groups of five (except that the last group may
     *  have fewer letters). */
    private void printMessageLine(String msg) {
        String groupOfFive = "";
        String finalOutput = "";
        for (int i = 0; i < msg.length(); i++) {
            groupOfFive += msg.charAt(i);
            if ((i + 1) % 5 == 0) {
                finalOutput += groupOfFive + " ";
                groupOfFive = "";
            } else if (i + 5 > msg.length()) {
                finalOutput += groupOfFive;
                groupOfFive = "";
            }
        }
        _output.println(finalOutput);
    }

    /** Alphabet used in this machine. */
    private Alphabet _alphabet;

    /** Number of rotors in machine. */
    private int _numRotors;

    /** Number of pawls in machine. */
    private int _pawls;

    /** An ArrayList of all available rotors. */
    private ArrayList<Rotor> _allRotors = new ArrayList<>();

    /** Source of input messages. */
    private Scanner _input;

    /** Source of machine configuration. */
    private Scanner _config;

    /** File for encoded/decoded messages. */
    private PrintStream _output;

    /** True if --verbose specified. */
    private static boolean _verbose;
}
