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


   POSIBILITY TO TEST TOPICS:
      //test : subscribe {"topics": ["usersTopic","tweetsTopic"]}
      //test:  unsubscribe {"topics": ["usersTopic"]}
      //again subscribe to users topic: subscribe {"topics": ["usersTopic"]}
      //again: subscribe {"topics": ["tweetsTopic"]}
