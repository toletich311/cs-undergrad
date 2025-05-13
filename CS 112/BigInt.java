/* 
 * BigInt.java
 *
 * A class for objects that represent non-negative integers of 
 * up to 20 digits.
 * 
 * name: Chandini Toleti
 * BUid: U29391556
 */

public class BigInt  {
    // the maximum number of digits in a BigInt -- and thus the length
    // of the digits array
    private static final int MAX_SIZE = 20;
    
    // the array of digits for this BigInt object
    private int[] digits;
    
    // the number of significant digits in this BigInt object
    private int sigDigits;
    private boolean overflow;

    /*
     * Default, no-argument constructor -- creates a BigInt that 
     * represents the number 0.
     */
    public BigInt() {
        digits = new int[MAX_SIZE];
        sigDigits = 1;  // 0 has one sig. digit--the rightmost 0!
	    overflow = false;
    }

    //constructor with array passed through 
    public BigInt(int[] arr){
        if(arr==null ||arr.length>MAX_SIZE || !(isValidArr(arr))){
            throw new IllegalArgumentException();
        } 
        digits = new int[MAX_SIZE];
        sigDigits=0; 
        for(int i =arr.length-1; i>=0; i--){
            digits[MAX_SIZE-sigDigits-1]= arr[i];
            sigDigits++; 
        }
        overflow = false; 
        sigDigits= findSig(digits);
    }

    //another constructor, but with passed int not arr
    public BigInt(int n){
        if(n<0){
            throw new IllegalArgumentException();
        }
        sigDigits= findSigInt(n);
        digits = new int[MAX_SIZE];

        for(int i = MAX_SIZE- 1; n>0; i--){
            digits[i]= n%10;
            n=n/10;
        }
    }

    //private helper, fings sig figures in an int, rather than an array
    private int findSigInt(int n){
        for(int i=0; i<MAX_SIZE;i++){
            if(n/Math.pow(10,i)<10){
                return i+1;

            }
        }
        return -1; 
    }




    //private helper method, determines first occurance of int n in array of ints arr
    private int findSig(int[] arr){
        //don't need to check if array is okay since it cannot get called if array is not okay!
        for(int i=0; i<MAX_SIZE;i++){
            if(arr[i]!=0){
                return MAX_SIZE-i; 
            }
        }
        return 1; 
    }
    
    //private helper method, checks if arr is filled with valid digits
    private boolean isValidArr(int[] arr){
        for(int i =0; i<arr.length; i++){
            if(arr[i]<0 || arr[i]>9){
                return false; 
            }
        }
        return true; 
    }


    //compares two BigInt objects , -1 if this<other, 0 if this==other, 1 if this>other
    public int compareTo(BigInt other){
        if(this.sigDigits<other.sigDigits){
            return -1;  
        } else if(this.sigDigits>other.sigDigits){
            return 1;
        }
        for(int i=MAX_SIZE-sigDigits; i<MAX_SIZE; i++){
            if(this.digits[i]==other.digits[i]){
                continue; 
            } else if(this.digits[i]>other.digits[i]){
                return 1; 
            } else{
                return -1; 
            }
        }
        return 0; 
    }

    
    public boolean greater_than_or_equal( BigInt other ){
        if (this.compareTo(other)>=0){
            return true; 
        }
        return false; 

    }

    public boolean smaller_than_or_equal( BigInt other ){
        if (this.compareTo(other)<=0){
            return true; 
        }
        return false; 
    }

    public boolean isOverFlow(){
        if (isOverFlow()){
            return true; 
        }
        return false; 
    }

    public BigInt add(BigInt other){
        if(other ==null){
            throw new IllegalArgumentException();
        }

        int carry=0; 
        int[] newA = new int[MAX_SIZE]; 

        for(int i=MAX_SIZE-1; i>=0; i--){
            if(this.digits[i]+other.digits[i]+carry>=10){
                newA[i]= (this.digits[i]+other.digits[i]+carry)%10; 
                carry = (this.digits[i]+other.digits[i]+carry)/10;
                //System.out.println(newA[i]); -->debugging
            } else{
                newA[i]= this.digits[i]+other.digits[i]+carry; 
                carry=0; 
                //System.out.println(newA[i]); -->debugging
            }
        
        }

        if (carry!=0){
            overflow=true; 
        }
        //newA = reverseArr(newA, findSig(newA));



        BigInt newInt = new BigInt(newA);
        return newInt; 
    }

  

    public BigInt mul(BigInt other){

        int num1=0; 
        int num2=0; 
        num1 = change(this.digits);
        num2 = change(other.digits);

        BigInt newInt = new BigInt(num1*num2); 
        return newInt; 
        /* 
        if(other == null){
            throw new IllegalArgumentException();
        }
        int[] add1 = new int[MAX_SIZE];
        int[] add2; 

        int carry =1; 
        for(int i = MAX_SIZE-1; i>=MAX_SIZE - this.sigDigits; i--){
            for(int j= MAX_SIZE-1; j>=MAX_SIZE - other.sigDigits; j--){
                if(this.digits[i]*other.digits[j]*carry>9){
                    add1[i]= (this.digits[i]*other.digits[j])%10;
                    carry = (this.digits[i]*other.digits[j])/10;
                    continue; 
                }
                add1[i]= (this.digits[i]*other.digits[j]*carry)%10;
                carry=1; 
            }
            add2 = arr2(add1, findSig(add1));
        }

        int[] newA= new int[MAX_SIZE]; 
        newA = arr
        */


    }


    //private helper method to change from array back to int
    private int change(int[] arr){
        int newNum=0; 
        int sigFig = findSig(arr); 
        for(int i=MAX_SIZE-sigFig-1; i<MAX_SIZE; i++){
            newNum += arr[i]*(Math.pow(10,MAX_SIZE-i-1)); 
        }
        return newNum; 
    }

    //private helper method to copy array elements
    /* 
    private int[] arr2(int[] arr1, int sigFig){
        int[] arr2 = new int[MAX_SIZE];
        for(int i =MAX_SIZE-sigFig;i<MAX_SIZE; i++ ){
            arr2[i] = arr1[i]; 
        }
        return arr2; 

    }
    */
 


    //accessor method , returns the value of the sigDigits field
    public int getNumSigDigits(){
        return this.sigDigits; 

    }

    public int getSigDigits(){
        return this.sigDigits; 

    }

    //accessor method -normally bad practice
    public int[] getDigits(){
        return this.digits; //returns reference to memory location, generally not good practice since this means we can mutate array from outside class
    }

    public String toString(){
        String str =""; 
        for(int i = MAX_SIZE- sigDigits; i<MAX_SIZE; i++){
            str+= digits[i]; 
        }
        return str; 
    }
    public static void main(String [] args) {
        System.out.println("Unit tests for the BigInt class.");
        System.out.println();

        /* 
         * You should uncomment and run each test--one at a time--
         * after you build the corresponding methods of the class.
         */

	
        System.out.println("Test 1: result should be 7");
        int[] a1 = { 1,2,3,4,5,6,7 };
        BigInt b1 = new BigInt(a1);
        System.out.println(b1.getNumSigDigits());
        System.out.println();
        
        System.out.println("Test 2: result should be 1234567");
        b1 = new BigInt(a1);
        System.out.println(b1);
        System.out.println();
        
        System.out.println("Test 3: result should be 0");
        int[] a2 = { 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 };
        BigInt b2 = new BigInt(a2);
        System.out.println(b2);
        System.out.println();

        
        
        System.out.println("Test 4: should throw an IllegalArgumentException");
        try {
            int[] a3 = { 0,0,0,0,23,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 };
            BigInt b3 = new BigInt(a3);
            System.out.println("Test failed.");
        } catch (IllegalArgumentException e) {
            System.out.println("Test passed.");
        } catch (Exception e) {
            System.out.println("Test failed: threw wrong type of exception.");
        }
        System.out.println();

        System.out.println("Test 5: result should be 1234567");
        b1 = new BigInt(1234567);
        System.out.println(b1);
        System.out.println();

        System.out.println("Test 6: result should be 0");
        b2 = new BigInt(0);
        System.out.println(b2);
        System.out.println();

        System.out.println("Test 7: should throw an IllegalArgumentException");
        try {
            BigInt b3 = new BigInt(-4);
            System.out.println("Test failed.");
        } catch (IllegalArgumentException e) {
            System.out.println("Test passed.");
        } catch (Exception e) {
            System.out.println("Test failed: threw wrong type of exception.");
        }
        System.out.println();


        
        System.out.println("Test 8: result should be 0");
        b1 = new BigInt(12375);
        b2 = new BigInt(12375);
        System.out.println(b1.compareTo(b2));
        System.out.println();

        System.out.println("Test 9: result should be -1");
        b2 = new BigInt(12378);
        System.out.println(b1.compareTo(b2));
        System.out.println();

        System.out.println("Test 10: result should be 1");
        System.out.println(b2.compareTo(b1));
        System.out.println();

        System.out.println("Test 11: result should be 0");
        b1 = new BigInt(0);
        b2 = new BigInt(0);
        System.out.println(b1.compareTo(b2));
        System.out.println();

        

        System.out.println("Test 12: result should be\n123456789123456789");
        int[] a4 = { 3,6,1,8,2,7,3,6,0,3,6,1,8,2,7,3,6 };
        int[] a5 = { 8,7,2,7,4,0,5,3,0,8,7,2,7,4,0,5,3 };
        BigInt b4 = new BigInt(a4);
        BigInt b5 = new BigInt(a5);
        BigInt sum = b4.add(b5);
        System.out.println(sum);
        System.out.println();

        System.out.println("Test 13: result should be\n123456789123456789");
        System.out.println(b5.add(b4));
        System.out.println();

        System.out.println("Test 14: result should be\n3141592653598");
        b1 = new BigInt(0);
        int[] a6 = { 3,1,4,1,5,9,2,6,5,3,5,9,8 };
        b2 = new BigInt(a6);
        System.out.println(b1.add(b2));
        System.out.println();

        System.out.println("Test 15: result should be\n10000000000000000000");
        int[] a19 = { 9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9,9 };    // 19 nines!
        b1 = new BigInt(a19);
        b2 = new BigInt(1);
        System.out.println(b1.add(b2));
        System.out.println();


        System.out.println("Test 16, mupltiplier: should be 255553");
        BigInt val1 = new BigInt(11111);
        BigInt val2 = new BigInt(23);
        BigInt product = val1.mul(val2);
        System.out.println(product);
    }
}
