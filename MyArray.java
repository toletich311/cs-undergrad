/* File: MyArrays
 *
 * Author: Chandini Toleti  CS112 
 *
 * Purpose: To create a class that allows you to
 * manipulate an array of integers.
 */

import java.util.Arrays;
import java.util.Scanner;                

public class MyArray  {

    private int SENTINEL = -999;// the sentinel value used to indicate end of input, initialized to -999
    private int DEFAULT_SIZE= 20;// the default size of the array if one is not specified, initialized to 20
    private int LOWER = 10;// the lower bound of the range of integer elements, initialized to 10
    private int UPPER = 50;// the upper bound of the range of integer elements, initialized to 50


    private int[] arr;// a data member to reference an array of integers
    private int numElements; // a data member to represent the number of elements entered into the array

    private int sum; //data member, repesents sum of all elements in the array
    private int min; //data member, represents the minimum value of all the elements in the array
    private int max; //data member, represents the maximum value of all the elemnts
    private double avg; //data member, average of all elements in the array

// CONSTRUCTORS
    // Initializes a MyArray object using default members
    public MyArray() {
       arr = new int[DEFAULT_SIZE];
       numElements = 0;
    }

    //4-1
    //custom contructor that creates an array of ints based on input argument
    public MyArray( int n ) {

        if(n>DEFAULT_SIZE){
            throw new IllegalArgumentException(); 
        }

        arr = new int[n]; 
        numElements = 0; 

    }

    /*4-2 custom constructor that creates object array based on array passed through method*/
    public MyArray( int[] array ) {
        if (array.equals(null)){
            throw new IllegalArgumentException(); 
        }
 
        arr = new int[array.length]; 
        numElements = 0; 

        for (int i =0; i< array.length; i++){
            if (array[i]>=LOWER && array[i]<=UPPER){
                arr[numElements] = array[i]; 
                numElements++;
            } 
        }

        computeStatistics();
    }

    /*
     * method prompts user to enter elements within the bounds
     * use may enter up to as many ints as the array can hold or less
     */

    public void inputElements() {

        //create scanner object
        Scanner console = new Scanner(System.in);

        System.out.println("   Enter up to `20` integers between `10` and `100` inclusive. Enter -999` to end user input: ");
        int val;

        if(this.numElements<this.arr.length){ //only takes input when there is space

            do{ 
                val = console.nextInt();
                if(val>=this.LOWER && val<=this.UPPER){
                    this.arr[this.numElements]= val;  
                    this.numElements++;  
                }

            }while(val!= this.SENTINEL && this.numElements<this.DEFAULT_SIZE && this.numElements<=this.arr.length); 
        }

        this.computeStatistics(); //compute statistics w/ updated elements
        
        console.close();  //close scanner
    }


    //4-4) determines whether an interger passed is a valid input
    public boolean validInput(int num) {

        if (num>=this.LOWER && num<=this.UPPER){
            return true; 
        }

        return false; 
    }

    //4-5) represnting the array as a string
    public  String toString() {

        if (this.numElements == 0){
            return ("[]");
        } //empty string

        String s = "[" ; 
        for (int i=0; i<this.numElements-1; i++){
            s+= this.arr[i] + ","; 
        }
        return (s + this.arr[this.numElements-1]+ "]"); 
    }


    //compute Statistics 4-6)
    public void computeStatistics() {
        //computes minimum, maximum, sum, and average 

        this.min = this.arr[0]; 
        this.max = this.arr[0];
        this.avg = 0.0; 
        this.sum = 0; 

        if(this.numElements>0){

            //use same loop iteration to calc min and max
        for(int i=0; i<this.numElements; i++){
            if (this.arr[i]<this.min){
                this.min = arr[i];
            }

            if(this.arr[i]>this.max){
                this.max = arr[i]; 
            }

            this.sum += this.arr[i]; 
        } 
        this.avg = ( (double)this.sum/this.numElements);
    }
    }

    //4-7) last index --> returns last occurence of the number passed to the method 
    public int lastIndex(int n) {
        for(int i =this.numElements-1; i>=0; i--){
            if (this.arr[i]==n){
                return i; 
            }
        }
        return -1; 
    }


    //4-8) insert -> inserts a specified number at specified position, returns false if input cannot be inserted
    public boolean insert(int n, int position) {

        if(position>this.numElements || this.numElements>=20 ||this.arr.length==this.numElements ){
            return false; 
        } else {
            for (int i =position+1; i<this.numElements+1; i++){
                this.arr[i]= this.arr[i-1];
            }
            this.arr[position]= n; 
            this.computeStatistics();
            return true; 
        }
    }


    //4-9) removeSmallest, removes smallest element in array, if no element is removed, return -1
    public int removeSmallest() {
        int minIdx=0; 
        this.computeStatistics(); 
        for(int i=0; i<this.numElements; i++){
            if (this.arr[i]==this.min){
                minIdx= i; 
            }
        }
        for(int j =minIdx; j<this.numElements-1; j++){
            this.arr[j] = this.arr[j+1]; 
        }
        return this.min; 

    }

    //4-10) grow --> grows list by n elements, returns false if list is too big to be grown by n elements
    public boolean grow(int n) {

        int newLen = n+this.numElements;
        if (newLen<=20){

            int[] newArr = new int[this.arr.length + n]; 
            for(int i=0; i<this.numElements; i++){
                newArr[i]= this.arr[i];
            }
            this.arr = newArr;
            this.computeStatistics();
            return true; 
        }
        return false; 

    }

    //4-11) shrink method --> shrinks array to be the exact number of elements, removes any empty spaces
    public boolean shrink() {

        if(this.arr.length == this.numElements || this.arr.length ==0 || this.numElements==0){
            return false; 
        } 
            int[] newArr= new int[(this.numElements)];

            for(int i=0; i<newArr.length; i++){
                newArr[i] = this.arr[i]; 
            }

            this.arr = newArr;
            this.computeStatistics();

  
            return true;

    }


    //accessor methods 
    public int getSum() {
        return this.sum; 
    }
    
    public int getMin() {
        return this.min; 
    }
    
    public int getMax() {
        return this.max; 
    }
    
    public double getAvg() {
        return this.avg; 
    }
    
    public int[] getArr() {
        return this.arr; 
    }
    
    public int getNumElements() {
        return this.numElements;  //return numElements count, different from size
    }
    
    public int getArrLength() {
        return this.arr.length; // return the current physical size (i.e. length) of the object's array
    }

    //4-13) computeHistogram --> turns array into a string histogram and returns it
    public String computeHistogram(){
        String s= ""; 
        for(int i=0; i<this.numElements; i++){
            for(int j=0; j<this.arr[i]; j++){
                s += "*"; 
            }
            s+= "\n"; 
        }
        return s; 
    }

    public static void main(String[] args) {
        
    

    

        System.out.println("\nUnit Test for MyArray.\n");
        //System.out.println(a2.computeHistogram());
        /* 
        System.out.println("array: " + arr);
        System.out.println("numElements: " +arr.numElements);
        System.out.println("length: "+ arr.getArrLength());

        System.out.println("array: " + a2.toString());
        System.out.println("numElements: " +a2.numElements);
        System.out.println("length: "+ a2.getArrLength());

        int[] arr = {10,15,20,25,30,35,40,50,0,0,0,0}; 
        MyArray a1 = new MyArray(arr);

        a1.shrink(); 
        //System.out.println(Arrays.toString(a1.arr));
        */

    }
}

