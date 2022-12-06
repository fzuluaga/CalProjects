package enigma;

import org.junit.Test;
import org.junit.Rule;
import org.junit.rules.Timeout;
import static org.junit.Assert.*;

import static enigma.TestUtils.*;

/** The suite of all JUnit tests for the Permutation class.
 *  @author Felipe Zuluaga
 */
public class PermutationTest {

    /** Testing time limit. */
    @Rule
    public Timeout globalTimeout = Timeout.seconds(5);

    /* ***** TESTING UTILITIES ***** */

    private Permutation perm;
    private String alpha = UPPER_STRING;

    /** Check that perm has an alphabet whose size is that of
     *  FROMALPHA and TOALPHA and that maps each character of
     *  FROMALPHA to the corresponding character of FROMALPHA, and
     *  vice-versa. TESTID is used in error messages. */
    private void checkPerm(String testId,
                           String fromAlpha, String toAlpha) {
        int N = fromAlpha.length();
        assertEquals(testId + " (wrong length)", N, perm.size());
        for (int i = 0; i < N; i += 1) {
            char c = fromAlpha.charAt(i), e = toAlpha.charAt(i);
            assertEquals(msg(testId, "wrong translation of '%c'", c),
                         e, perm.permute(c));
            assertEquals(msg(testId, "wrong inverse of '%c'", e),
                         c, perm.invert(e));
            int ci = alpha.indexOf(c), ei = alpha.indexOf(e);
            assertEquals(msg(testId, "wrong translation of %d", ci),
                         ei, perm.permute(ci));
            assertEquals(msg(testId, "wrong inverse of %d", ei),
                         ci, perm.invert(ei));
        }
    }

    /* ***** TESTS ***** */

    @Test
    public void checkIdTransform() {
        perm = new Permutation("", UPPER);
        checkPerm("identity", UPPER_STRING, UPPER_STRING);
    }

    @Test
    public void testSize() {
        Alphabet alpha1 = new Alphabet("ABCDEFG");
        assertEquals(7, alpha1.size());

        Alphabet alpha2 = new Alphabet();
        assertEquals(26, alpha2.size());

        Alphabet alphaSmall = new Alphabet("A");
        assertEquals(1, alphaSmall.size());
    }

    @Test
    public void testContains() {
        Alphabet alpha1 = new Alphabet("ABCDEFG");
        assertTrue(alpha1.contains('A'));
        assertTrue(alpha1.contains('G'));
        assertFalse(alpha1.contains('Z'));
        assertFalse(alpha1.contains('Q'));
    }

    @Test
    public void testToChar() {
        Alphabet alpha1 = new Alphabet("ABCDEFG");
        assertEquals('A', alpha1.toChar(0));
        assertEquals('G', alpha1.toChar(6));

        Alphabet alphaNormal = new Alphabet();
        assertEquals('A', alphaNormal.toChar(0));
        assertEquals('Z', alphaNormal.toChar(25));
    }

    @Test
    public void testToInt() {
        Alphabet alpha1 = new Alphabet("ABCDEFG");
        assertEquals(0, alpha1.toInt('A'));
        assertEquals(6, alpha1.toInt('G'));

        Alphabet alphaNormal = new Alphabet();
        assertEquals(25, alphaNormal.toInt('Z'));
        assertEquals(0, alphaNormal.toInt('A'));
    }

    @Test
    public void testPermute() {
        Permutation p = new Permutation("(BACD)",
                        new Alphabet("ABCDE"));
        assertEquals(2, p.permute(0));
        assertEquals(0, p.permute(1));
        assertEquals(3, p.permute(2));
        assertEquals(1, p.permute(3));
        assertEquals(4, p.permute(4));

        assertEquals('A', p.permute('B'));
        assertEquals('C', p.permute('A'));
        assertEquals('D', p.permute('C'));
        assertEquals('B', p.permute('D'));
        assertEquals('E', p.permute('E'));

        Permutation p1 = new Permutation("(BA) (ED) (FC)",
                         new Alphabet("ABCDEFG"));
        assertEquals('A', p1.permute('B'));
        assertEquals('B', p1.permute('A'));
        assertEquals('D', p1.permute('E'));

        assertEquals(0, p1.permute(1));
        assertEquals(2, p1.permute(5));
        assertEquals(6, p1.permute(6));
    }

    @Test
    public void testInvert() {
        Permutation p = new Permutation("(BACD)",
                        new Alphabet("ABCDE"));
        assertEquals(1, p.invert(0));
        assertEquals(3, p.invert(1));
        assertEquals(0, p.invert(2));
        assertEquals(2, p.invert(3));
        assertEquals(4, p.invert(4));

        assertEquals('B', p.invert('A'));
        assertEquals('A', p.invert('C'));
        assertEquals('C', p.invert('D'));
        assertEquals('D', p.invert('B'));
        assertEquals('E', p.invert('E'));

        Permutation p1 = new Permutation("(BA) (ED) (FC)",
                         new Alphabet("ABCDEFG"));
        assertEquals('A', p1.invert('B'));
        assertEquals('E', p1.invert('D'));
        assertEquals('G', p1.invert('G'));

        assertEquals(1, p1.invert(0));
        assertEquals(2, p1.invert(5));
        assertEquals(6, p1.invert(6));
    }

    @Test
    public void testOutOfIndex() {
        Permutation p = new Permutation("(ABCDEFGHIJKLMNOPQRSTUVWXYZ)",
                        new Alphabet());
        assertEquals(0, p.permute(51));
        assertEquals(13, p.permute(38));
        assertEquals(3, p.invert(30));
        assertEquals(23, p.invert(50));
    }

    @Test
    public void testAlphabet() {
        Alphabet alpha1 = new Alphabet();
        Permutation p = new Permutation("", alpha1);
        assertEquals(alpha1, p.alphabet());

        Alphabet alpha2 = new Alphabet("ABC");
        Permutation p1 = new Permutation("", alpha2);
        assertEquals(alpha2, p1.alphabet());

        Permutation p2 = new Permutation("(A)",
                         new Alphabet("A"));
        assertEquals('A', p2.permute('A'));
        assertEquals(0, p2.permute(0));
        assertEquals('A', p2.invert('A'));
        assertEquals(0, p2.invert(0));
    }

    @Test
    public void testDerangement() {
        Permutation p = new Permutation("(BACD)(E)",
                        new Alphabet("ABCDE"));
        assertFalse(p.derangement());

        Permutation p1 = new Permutation("(CBA)",
                         new Alphabet("ABCD"));
        assertFalse(p1.derangement());

        Permutation p2 = new Permutation("(ABCD)",
                         new Alphabet("ABCD"));
        assertTrue(p2.derangement());

        Permutation p3 = new Permutation("(AB)",
                         new Alphabet("AB"));
        assertTrue(p3.derangement());

        Permutation p4 = new Permutation("(AB)(CD)(EFG)",
                         new Alphabet("ABCDEFG"));
        assertTrue(p4.derangement());

        Permutation p5 = new Permutation("(A)(B)(C)(D)",
                         new Alphabet("ABCD"));
        assertFalse(p5.derangement());

        Permutation p6 = new Permutation("(A)",
                         new Alphabet("A"));
        assertFalse(p6.derangement());
    }

    @Test(expected = EnigmaException.class)
    public void testNotInAlphabet() {
        Permutation p = new Permutation("(BACD)",
                        new Alphabet("ABCD"));
        p.invert('F');
        p.invert(7);
        p.permute('F');
        p.invert(7);
    }

}
