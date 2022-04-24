import Head from "next/head";
import { useState } from "react";
import { createActivity } from "../service";
import Datetime from 'react-datetime';
import "react-datetime/css/react-datetime.css";

export default function Home(props) {
  /** You can access to liff and liffError object through the props.
   *  const { liff, liffError } = props;
   *  console.log(liff.getVersion());
   *
   *  Learn more about LIFF API documentation (https://developers.line.biz/en/reference/liff)
   **/
  const { liff, liffError } = props;
  const [activity, setActivity] = useState();
  const [date, setDate] = useState(new Date());
  const [place, setPlace] = useState();

  const handleSubmit = async (e) => {
    // Post to golang server
    e.preventDefault()
    try {
      const idToken = liff.getIDToken()
      await createActivity({ activity, date, place, idToken });

      // Submit message
      
      await liff.sendMessages([
        {
          "type": "text",
          "text": "我要舉辦活動!"
        }
      ]);
      console.log("Success!!")
    } catch (err) {
      console.log("err");
      console.log(err);
    }

    console.log("Success22!!")
  };

  return (
    <div>
      <Head>
        <title>Activity Scheduler</title>
      </Head>
      <div className="home">
        <h1 className="home__title">Fill in your activity!</h1>
        <form onSubmit={handleSubmit} method="post">
          <label id="activity">Activity Name:</label>
          <input
            type="text"
            id="activity"
            name="activity"
            required
            onChange={(e) => setActivity(e.target.value)}
          />
          <br />
          <Datetime />
          <br />
          <label id="name">Place:</label>
          <input
            type="text"
            id="place"
            name="place"
            required
            onChange={(e) => setPlace(e.target.value)}
          />

          <button type="submit">Submit</button>
        </form>
      </div>
    </div>
  );
}
