# tweeter-sentiment-analyzer

1. Link to video with output result of first laboratory work:

    https://drive.google.com/file/d/1-f9nar_ijPz-_hL8d3zMtA0kS_fXVpix/view?usp=sharing



2. Link to video with output result of second laboratory work:

   https://drive.google.com/drive/folders/1jlEr2jJOYAcp1Nka_7YiUOerO-rKfH-E?usp=sharing

3. Link to video with output result of third laboratory work:
   https://drive.google.com/file/d/1y9OhHsrlTlvvpokxwF5qBLqTliDXncOy/view?usp=sharing

        
    DOCKER STUFF:
        1. sudo docker-compose build --no-cache
        2. sudo docker-compose up
        3. Connects to Broker using telnet -> telnet 127.0.0.1 8088;


    TO TEST TOPICS USING DURABLE QUEUES:
        

    SUBSCRIBE:

        {"topics": [{"value": "tweetsTopic","is_durable": true},{"value": "usersTopic","is_durable": false}],"command":"subscribe"}
        
        {"topics": [{"value": "usersTopic","is_durable": false}],"command":"subscribe"}
        
        {"topics": [{"value": "tweetsTopic","is_durable": true}],"command":"subscribe"}

    UNSUBSCRIBE:
        {"topics": [{"value": "usersTopic"}],"command":"unsubscribe"} // from non durable
        
        {"topics": [{"value": "tweetsTopic"}],"command":"unsubscribe"} //from durable topic -> nothing happens;

    STOP COMMAND:
         {"command":"stop"}
   
    UNiQUE ID MSG:
        ***** value => value is the string unique id which can be taken from logs; *******
        {"unique_id_for_durable": "value"} 

