package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
	"os"
	"encoding/json"
    "encoding/binary"
    "encoding/hex"
    "strconv"
    "database/sql"
    _ "mysql-driver/mysql-master"
	"time"
    )

func main() {
    db, err := sql.Open("mysql", "root:@/test")
	  if err != nil {
        panic(err.Error())
    }
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
    var result map[string]interface{}
    var ver string
	var h int
    // if there is an error opening the connection, handle it
	//row := db.QueryRow("SELECT MAX(height) + 1 as next_number FROM hashes")
    //err = row.Scan(&h)
    //if err != nil {
    //if err == sql.ErrNoRows {
     //   fmt.Println("Zero rows found")
      //  } else {
      //  panic(err)
     //   }
    //}
	sql := "SELECT h0.height FROM hashes h0 " +
	" JOIN hashes h1 ON h1.height = h0.height + 1 "  +
	" WHERE HEX(h1.previous_block) <> SHA2( UNHEX (SHA2(CONCAT(h0.VERSION, h0.previous_block, h0.merkle_root, h0.TIME, h0.bits, h0.nonce),256) ),256)"
    results, err := db.Query(sql)
    if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}	
    for results.Next() {
	 err = results.Scan(&h)
	 if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
    // for h<556000 {
    
      timeout := time.Duration(11 * time.Second)
	  client := http.Client{
        Timeout: timeout,
      }
      response, err := client.Get("https://blockchain.info/block-height/"+ strconv.Itoa(h)+"?format=json")
      if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
		} else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
              
	       if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
		 }
		json.Unmarshal([]byte(contents), &result)
		block := result["blocks"].([]interface{})
		//hash := block[0].(map[string]interface{})["hash"]
		prev_block := reverse(block[0].(map[string]interface{})["prev_block"].(string))
        mrkl_root := reverse(block[0].(map[string]interface{})["mrkl_root"].(string))
       /* if block[0].(map[string]interface{})["ver"].(float64) == 1 {
          ver = "01000000"
        } else  {
          ver = "02000000"
        }*/
		version := uint32(block[0].(map[string]interface{})["ver"].(float64));
		nonce :=uint32(block[0].(map[string]interface{})["nonce"].(float64))
        buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, version)
		ver = hex.EncodeToString(buf)
        binary.LittleEndian.PutUint32(buf, nonce)
        encoded_nonce := hex.EncodeToString(buf)
        //nonce = strconv.FormatFloat(block[0].(map[string]interface{})["nonce"].(float64), 'f', -1, 64)
		time :=uint32(block[0].(map[string]interface{})["time"].(float64))
        binary.LittleEndian.PutUint32(buf, time)
        encoded_time := hex.EncodeToString(buf)
        bits :=uint32(block[0].(map[string]interface{})["bits"].(float64))
        binary.LittleEndian.PutUint32(buf, bits)
        encoded_bits := hex.EncodeToString(buf)
        //query := "insert into hashes(height,version, previous_block, merkle_root,time, bits,nonce) values ("+ strconv.Itoa(h) + ",0x" + ver + ", 0x" + prev_block + 
        //", 0x" + mrkl_root +", 0x" + encoded_time + ", 0x" + encoded_bits + ", 0x" + encoded_nonce + ") "
		query := "update hashes set version=0x" + ver + ", previous_block = 0x"+ prev_block + ", merkle_root = 0x" + mrkl_root + ",time = 0x"+ encoded_time + ", bits = 0x" + encoded_bits + ", nonce = 0x" + encoded_nonce + 
		" where height  =  "+ strconv.Itoa(h)

		
       // fmt.Printf("%s\n", string(contents))
		fmt.Printf("%s\n", query)

        rows, err := db.Query(query)
         if err != nil {
           panic(err.Error()) // proper error handling instead of panic in your app
         } else
         {
             fmt.Println(rows)
         }
        //fmt.Println(hash, prev_block)
		//fmt.Printf("%d\n", int64(nonce))
		
		rows.Close()
     }
	 time.Sleep(1300 * time.Millisecond)
	 if h%995 ==0  {
	  time.Sleep(10000 * time.Millisecond) 
	  }
	 //h = h+1
	}
   }
  func reverse(s string) string {
    runes := []rune(s)
     for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
      runes[i], runes[j] = runes[j], runes[i]
     }
	 for i:= 0; i < len(runes); i = i+2 {
	   runes[i],runes[i+1] = runes[i+1], runes[i]
	  }
    return string(runes)
  }
  
 
