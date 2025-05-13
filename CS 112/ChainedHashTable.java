/*
 * ChainedHashTable.java
 *
 * Computer Science 112, Boston University
 * 
 * Modifications and additions by:
 *     name: Chandini Toleti
 *     email: toletich@bu.edu
 */

import java.util.*;     // to allow for the use of Arrays.toString() in testing

/*
 * A class that implements a hash table using separate chaining.
 */
public class ChainedHashTable implements HashTable {
    /* 
     * Private inner class for a node in a linked list
     * for a given position of the hash table
     */
    private class Node {
        private Object key;
        private LLQueue<Object> values;
        private Node next;
        
        private Node(Object key, Object value) {
            this.key = key;
            values = new LLQueue<Object>();
            values.insert(value);
            next = null;
        }
    }
    
    private Node[] table;      // the hash table itself
    private int numKeys;       // the total number of keys in the table
        
    /* hash function */
    public int h1(Object key) {
        int h1 = key.hashCode() % table.length;
        if (h1 < 0) {
            h1 += table.length;
        }
        return h1;
    }
    
    /*** Add your constructor here ***/
    public ChainedHashTable(int size){

        //check if valid size
        if(size<0){
            throw new IllegalArgumentException();
        }

        //intialize both fields 
        table = new Node[size]; 
        numKeys = 0; 

    }

    //accessor methods
    public int getNumKeys(){
        return this.numKeys; 
    }

    public double load(){
        return((double)numKeys/table.length);
    }
    
    public Object[] getAllKeys(){
        Object[] keys = new Object[numKeys];
        int count = numKeys; 
        for(int i=0; i<table.length; i++){
            if(table[i]!=null){
                Node travNode = table[i];
                while(travNode!=null){
                    keys[numKeys-count]= travNode.key;
                    travNode = travNode.next;
                    count--; 
                
                }
            }
        } 
        return keys; 
    }

    public void resize(int newSize){
        if(newSize<0 || newSize == table.length){
            throw new IllegalArgumentException();
        }

        Node[] newTable = new Node[newSize]; 


        Object[] keys = new Object[numKeys];
        keys = this.getAllKeys();

        this.table = newTable; 

        System.out.println(Arrays.toString(keys));

        for(int i=0; i<keys.length; i++){
            this.insert(keys[i], keys[i]);
        }
        numKeys-= keys.length; 
    }
    
    /*
     * insert - insert the specified (key, value) pair in the hash table.
     * Returns true if the pair can be added and false if there is overflow.
     */
    public boolean insert(Object key, Object value) {
        /** Replace the following line with your implementation. **/
        if (key == null) {
            throw new IllegalArgumentException("key must be non-null");
        }

        int idx = h1(key); 
        //System.out.println(idx + "is new Val!!"); --> debugging, ignore
        Node newNode = new Node(key, value);
        if (table[idx]== null){ 
            table[idx] = newNode;

        }else{ //has to be more effecient way 

            //check for dublicate 

            while(table[idx].next!=null){
                if (table[idx].next.key.equals(key)){
                    return false; //return false --> don't insert or update numKeys
                }
                table[idx] = table[idx].next;
            }


            Node shit = table[idx];
            newNode.next= shit; 
            table[idx] = newNode;

            //original reverse order
            //table[idx].values.insert(value);
        }
 
        this.numKeys++; 
        return true;
    }
    
    /*
     * search - search for the specified key and return the
     * associated collection of values, or null if the key 
     * is not in the table
     */
    public Queue<Object> search(Object key) {
        /** Replace the following line with your implementation. **/
      
            if (key == null) {
                throw new IllegalArgumentException("key must be non-null");
            }
            
            int i = h1(key);
            
            if (i == -1 || table[i] == null) {
                return null;
            } else {
                return table[i].values;
            }
        
   
    }
    
    /* 
     * remove - remove from the table the entry for the specified key
     * and return the associated collection of values, or null if the key 
     * is not in the table
     */
    public Queue<Object> remove(Object key) {
        /** Replace the following line with your implementation. **/
        if (key == null) {
            throw new IllegalArgumentException("key must be non-null");
        }

        int i = h1(key);

        if (i == -1 || table[i] == null) {
            return null;
        }
        int count=0; 

        Node trav = table[i];
        while(trav.next!=null){
            trav = trav.next;
            count++; 
        }
    
        
        LLQueue<Object> removedVals = trav.values;
        trav=null;

        table[i] = null;
        if(count>0){
            numKeys-=count;
        }
   
        numKeys--;
        return removedVals;
        

   
    }
    
   
    
    
    
    /*
     * toString - returns a string representation of this ChainedHashTable
     * object. *** You should NOT change this method. ***
     */
    public String toString() {
        String s = "[";
        
        for (int i = 0; i < table.length; i++) {
            if (table[i] == null) {
                s += "null";
            } else {
                String keys = "{";
                Node trav = table[i];
                while (trav != null) {
                    keys += trav.key;
                    if (trav.next != null) {
                        keys += "; ";
                    }
                    trav = trav.next;
                }
                keys += "}";
                s += keys;
            }
        
            if (i < table.length - 1) {
                s += ", ";
            }
        }       
        
        s += "]";
        return s;
    }

    public static void main(String[] args) {

        /** Add your unit tests here **/
        ChainedHashTable table = new ChainedHashTable(5);
        table.insert("howdy", 15);
        table.insert("goodbye", 10);
        System.out.println(table.insert("apple", 5));
        System.out.println(table);

        //test remove

        table.remove("howdy");
        System.out.println(table);

        //test search 
        System.out.println(table.search("goodbye"));
        System.out.println(table.search("howdy"));
        System.out.println(table.search("apple"));

        //test getAll keys
        Object[] keys = table.getAllKeys();
        System.out.println(Arrays.toString(keys));


        //test numKeys
        System.out.println(table.getNumKeys());

        //load test
        table.insert("pear", 6);
        System.out.println(table.load());

        //resize test 
        table.resize(7);
        System.out.println(table);




   }
}
