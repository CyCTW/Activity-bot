import Head from "next/head";
import Link from 'next/link'

import { useState } from "react";
import { createActivity } from "../service";
import Datetime from "react-datetime";
import "react-datetime/css/react-datetime.css";
import {
  Button,
  Center,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
// import DateTimePicker from 'react-datetime-picker';

export default function Home(props) {
  /** You can access to liff and liffError object through the props.
   *  const { liff, liffError } = props;
   *  console.log(liff.getVersion());
   *
   *  Learn more about LIFF API documentation (https://developers.line.biz/en/reference/liff)
   **/
  const { liff, liffError } = props;
  const [name, setName] = useState();
  const [date, setDate] = useState();
  const [place, setPlace] = useState();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    // Post to golang server
    e.preventDefault();
    setLoading(true);
    try {
      const idToken = liff.getIDToken();
      const dateString = date.toDate().toJSON()
      const res = await createActivity({ name, dateString, place, idToken });
      const activityID = res.data?.activityID;
      // Submit message

      await liff.sendMessages([
        {
          type: "text",
          text: `我要舉辦活動, ID: ${activityID}`,
        },
      ]);
      console.log("Success!!");
    } catch (err) {
      console.log("err");
      console.log(err);
    }
    setLoading(false);

    console.log("Success22!!");
  };

  const handleDate = (date) => {
    // const jj = date.toDate();
    // console.log(jj.toJSON());
    // console.log(typeof jj);
    setDate(date);
  };
  return (
    <div
    style={{
      backgroundColor: 'gray',
    }}
    >
      <Head>
        <title>Activity Scheduler</title>
      </Head>
      <Link href="/activity/second"><a>Second Post</a></Link>
      <div className="home">
        <h1 className="home__title">Fill in your activity!</h1>
        <FormControl>
          <VStack mt={10}>
            <FormLabel htmlFor="activity">活動名稱</FormLabel>
            <input
              type="text"
              id="activity"
              name="activity"
              required
              width="auto"
              onChange={(e) => setName(e.target.value)}
            />
            <FormLabel htmlFor="date">日期</FormLabel>
            <Datetime onChange={handleDate} value={date} />
            <FormLabel htmlFor="place">地點</FormLabel>
            <input
              type="text"
              id="place"
              name="place"
              required
              width="auto"
              onChange={(e) => setPlace(e.target.value)}
            />
            <Button mt={10} colorScheme='blue' type="submit" isLoading={loading} onClick={handleSubmit}>
              Submit
            </Button>
          </VStack>
        </FormControl>
        {/* <form onSubmit={handleSubmit} method="post">
          <label id="activity">Activity Name:</label>
          <input
            type="text"
            id="activity"
            name="activity"
            required
            onChange={(e) => setName(e.target.value)}
          />
          <br />
          <label id="date">Date:</label>
          <input
            type="text"
            id="date"
            name="date"
            required
            onChange={(e) => setDate(e.target.value)}
          />
          <Datetime onChange={handleDate} value={date} />
          <br />
          <label id="name">Place:</label>
          <input
            type="text"
            id="place"
            name="place"
            required
            onChange={(e) => setPlace(e.target.value)}
          />

          <Button type="submit" isLoading={loading}>Submit</Button>
        </form> */}
      </div>
    </div>
  );
}
