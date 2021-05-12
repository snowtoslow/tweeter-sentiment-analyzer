# tweeter-sentiment-analyzer

1. Link to video with output result of first laboratory work:

    https://drive.google.com/file/d/1-f9nar_ijPz-_hL8d3zMtA0kS_fXVpix/view?usp=sharing



2. Link to video with output result of second laboratory work:

   https://drive.google.com/drive/folders/1jlEr2jJOYAcp1Nka_7YiUOerO-rKfH-E?usp=sharing


3.
   DOCKER STUFF:
      1. sudo docker-compose build --no-cache
      2. sudo docker-compose up


   RUN TWEETS SERVER:
         sudo docker run -p 4000:4000 alexburlacu/rtp-server:faf18x
         telnet 127.0.0.1 8088


   POSSIBILITY TO TEST TOPICS BEFORE ADDING DURABLE QUEUES:
      //test : subscribe {"topics": ["usersTopic","tweetsTopic"]}
      //test:  unsubscribe {"topics": ["usersTopic"]}
      //again subscribe to users topic: subscribe {"topics": ["usersTopic"]}
      //again: subscribe {"topics": ["tweetsTopic"]}


   TO TEST TOPICS USING DURABLE QUEUES:
   {"topics": [{"value": "tweetsTopic","is_durable": true},{"value": "usersTopic","is_durable": false}],"command":"subscribe"}

   {"topics": [{"value": "tweetsTopic","is_durable": true}],"command":"subscribe"}

   {"command":"stop"}

   {"unique_id_for_durable":"2df37c44-3326-077f-960d-92c6222bb5b7"}

   {"topics": [{"value": "usersTopic","is_durable": false}],"command":"subscribe"}

   {"topics": [{"value": "tweetsTopic"}],"command":"unsubscribe"}
