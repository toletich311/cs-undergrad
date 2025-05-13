/*
 * Student First Name: Chandini
 * Student Last Name: Toleti
 * Student BU Number: U29391556
 * Purpose: creates a class of a set and methods to manipulate any set with other sets or compare with other sets. 
 */

import java.util.Arrays;

public class Set  {

    private static final int DEFAULT_SIZE = 10; // default size of initial set    
    private int[] set;      // array referece to the set
    private int next;       // index to next available slot in the set array
    
    /*
     * Constructors
     */
    public Set() {
        // your code here
        this.set = new int[DEFAULT_SIZE]; //?? or should i use default size
        next = 0; 
    }

    /*
     * Implement remaining methods as specified by the problem set
     */
 

    public Set(int[] arr){
        this.set = new int[arr.length]; 
        next=0; 
        boolean flag = true; 
        for(int i =0; i<arr.length; i++){
            for(int j=0; j<this.set.length; j++){
                if(arr[i]==this.set[j]){
                    flag=false; 
                }
                
            }
            if(flag){
                this.set[next]=arr[i];
                next++; 
            }
            flag = true; 
        }

    }

    public Set clone(){
        Set copy = new Set(this.set);
        return copy; 
    }

    public String toString(){
        if(this.next==0){
            return "[]"; 
        }
        String s= "["; 
        for(int i=0; i<this.cardinality(); i++){
            if (i!= (this.cardinality()-1)){
                s+= this.set[i]+ ",";
                continue; 
            }
            s+= set[i];
        }
        return s+= "]";
    }

    public int cardinality(){
        return this.next; 
    }

    public  boolean isEmpty(){
        if(this.next==0){
            return true; 
        }
        return false; 
    }

    public boolean member(int n){
       // if(this.isEmpty()){
         //   return false; 
        //}
        for(int i=0; i<this.cardinality(); i++){
            if(this.set[i]==n){
                return true; 
            }
        }
        return false; 
    }
    public boolean subset(Set S){
        //if more unique elements cannot be subset 
        if(this.cardinality()>S.cardinality()){
            return false;
        }

        for(int i=0; i<this.cardinality();i++){
            if(S.member(this.set[i])==false){
                return false; 
            }
        }

        return true; 
    }

    public boolean equal(Set S){

        //if not same num elements or NOT subset, obviously not same elements
        if(this.cardinality()!=S.cardinality() || this.subset(S)==false){
            return false; 
        }

        //check if each element of S is also an element of this
        for(int i=0; i<this.cardinality(); i++){
            if (!(this.member(S.set[i]))){
                return false; 

            }
        }

        return true; 

    }

    public void insert(int k) {
        if(!(this.member(k))){
            if(this.next == this.set.length){
                this.grow(); 
            }
            this.set[next]= k; 
            next++; 

        }
    }

    public void delete(int k){
        int[] newA= new int[this.cardinality()];
        int j=0; 
        if(this.member(k)){
            if(this.cardinality()==1){
                this.set=new int[DEFAULT_SIZE]; 
                next =0; 
            }

            for(int i =0; i<this.cardinality(); i++){
                if(this.set[i]==k){
                    newA[i]=this.set[i+1];  
                }else{
                    newA[j]=this.set[i];
                    j++; 
                }
            }
            this.set = newA; 
            next = findCard(newA); 
        }

    }

    //helper method to find cardinality 
    private int findCard(int[] arr){
        int count=0; 
        for(int i=0; i<arr.length; i++){
            if(arr[i]!=0){
                count++; 
            }
        }
        return count; 

    }

    public void grow(){
        int[] big = new int[this.set.length+DEFAULT_SIZE];
        for(int i=0; i<this.set.length; i++){
            big[i]=this.set[i]; 
        }
        this.set = big; 
    }

    public Set union(Set S){
       
        
        int[] unionA = new int[this.cardinality()+S.cardinality()];


        for(int i=0; i<this.cardinality();i++){
            unionA[i]= this.set[i]; 
        }
        for(int j=this.cardinality(); j<this.cardinality()+S.cardinality(); j++){
            unionA[j]= S.set[j-this.cardinality()]; 
        }

        Set U1 = new Set(unionA);

        /* 
        for(int i =0; i<this.cardinality(); i++){
            U1.set[i]= this.set[i]; 
            //next++; 
        }

        for(int j=0; j<S.cardinality(); j++){
            if(!(U1.member(S.set[j]))){
                U1.insert(S.set[j]);
                //next++; 
            }
        }
        */

        return U1;
    }

    public Set intersection(Set S){
        Set I1 = new Set(new int[this.cardinality()+S.cardinality()]);

        for(int i=0; i<this.cardinality(); i++){
            if (S.member(this.set[i])){
                I1.insert(this.set[i]);
            }
        }
        return I1;
    }

    public Set setdifference(Set S){
        Set D1 = new Set(new int[this.cardinality()+S.cardinality()]);

        for(int i=0; i<this.cardinality(); i++){
            if (!(S.member(this.set[i]))){
                D1.insert(this.set[i]);
            }
        }
        return D1;
    }


    public static void main(String [] args){

    }
}