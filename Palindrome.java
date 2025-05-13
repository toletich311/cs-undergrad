/*
 * Palindrome.java
 *
 * Computer Science 112
 *
 * Modifications and additions by:
 *     name: Chandini Toleti
 *     username: toletich@bu.edu
 */
   
public class Palindrome {
    // Add your definition of isPal here.
    public static boolean isPal(String str){
        //check for illegal arguement
        if(str==null){
            throw new IllegalArgumentException();
        }

        String strT = str.toLowerCase(); 
        int len = str.length(); 
        
        if(len==0 ||len==1){
            return true;
        }

        LLStack<Character> forward = new LLStack<Character>();
        LLQueue<Character> backward = new LLQueue<Character>(); 

        for(int i=0; i<str.length(); i++){
            if(strT.charAt(i)>=97&&strT.charAt(i)<=122){
                forward.push(strT.charAt(i));
                backward.insert(strT.charAt(i));
            }

        }
        System.out.println(backward);
        System.out.println(forward);

        for(int j=0; j<len; j++){
            if(forward.pop() != backward.remove()){
                return false; 
            }
        }
        return true; 

    }
    
    public static void main(String[] args) {
        System.out.println("--- Testing method isPal ---");
        System.out.println();

        //test 0 ------
        System.out.println("(0) Testing on \"A man, a plan, a canal, Panama!\"");
        try {
            boolean results = isPal("A man, a plan, a canal, Panama!");
            boolean expected = true;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println();  // include a blank line between tests

        /*
         * Add five more unit tests that test a variety of different
         * cases. Follow the same format that we have used above.
         */

        //test 1 -----
        System.out.println("(1) Testing on \"Wow! this is NOT one!\"");
        try {
            boolean results = isPal("Wow! this is NOT one!");
            boolean expected = false;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println();

        //test 2 ---------
        System.out.println("(2) Testing on \"Racecar Bob, RADAR BOBRACE , car!\"");
        try {
            boolean results = isPal("Racecar Bob, RADAR BOBRACE , car!");
            boolean expected = true;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println();

        //test 3 -------
        System.out.println("(3) Testing on \"Woah there, Jamal!\"");
        try {
            boolean results = isPal("Woah there, Jamal!");
            boolean expected = false;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println();

        //test 4 -------
        System.out.println("(4) Testing on \"Madam, I'm Adam?\"");
        try {
            boolean results = isPal("Madam, I'm Adam?");
            boolean expected = true;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println(); 
        
        
        //test 5 ---------
        System.out.println("(5) Testing on \"PULL UP, if I pull up...\"");
        try {
            boolean results = isPal("PULL UP, if I pull up...");
            boolean expected = true;
            System.out.println("actual results:");
            System.out.println(results);
            System.out.println("expected results:");
            System.out.println(expected);
            System.out.print("MATCHES EXPECTED RESULTS?: ");
            System.out.println(results == expected);
        } catch (Exception e) {
            System.out.println("INCORRECTLY THREW AN EXCEPTION: " + e);
        }
        System.out.println();
        
        
    }
}