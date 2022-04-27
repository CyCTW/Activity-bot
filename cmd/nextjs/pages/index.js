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
  useToast,
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
  const toast = useToast()

  const handleSubmit = async (e) => {
    // Post to golang server
    e.preventDefault();
    setLoading(true);
    try {
      const idToken = liff.getIDToken();
      const dateString = date.toDate().toJSON()
      const res = await createActivity({ name, dateString, place, idToken });
      const activityName = res.data?.activityName;
      // Submit message

      await liff.sendMessages([
        {
          type: "text",
          text: `我舉辦了活動 ${activityName} !`,
        },
        {
          type: "text",
          text: `@顯示活動-${activityName}`,
        },
      ]);
      console.log("Success!!");
    } catch (err) {
      console.log("err");
      console.log(err);
    }
    setLoading(false);

    toast({
      title: 'Activity created.',
      description: "We've created activity for you.",
      position: 'bottom',
      status: 'success',
      duration: 5000,
      isClosable: true,
    })
  };

  const handleDate = (date) => {
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
      <div className="home">
        <h1 className="home__title">填寫活動!</h1>
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
      </div>
    </div>
  );
}
