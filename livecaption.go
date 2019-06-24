// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command livecaption pipes the stdin audio data to
// Google Speech API and outputs the transcript.
//
// As an example, gst-launch can be used to capture the mic input:
//
//    $ gst-launch-1.0 -v pulsesrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption
package main

// [START speech_transcribe_streaming_mic]
import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

	"io/ioutil"


	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"


)

// func main() {
func liveCaption() {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Send the initial configuration message.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		// Pipe stdin to the API.
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {

				// function doSomething() {
				// 	if (error) return error
				// 	console.log('no error!')
				// 	return null
				// }

				// const err = doSomething()
				// if (err) {
				// 	handleError()
				// }

				// if err = doSomething(); err != nil {
				// 	handleError()
				// }

				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: buf[:n],
					},
				}); err != nil {
					log.Printf("Could not send audio: %v", err)
				}
			}
			// end-of-file error
			if err == io.EOF {
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Could not close stream: %v", err)
				}
				return
			}
			if err != nil {
				log.Printf("Could not read from stdin: %v", err)
				continue
			}
		}
	}()

	for {
		// function returnNumberPair() {
		// 	return [5, 10]
		// }
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot stream results: %v", err)
			// log.Fatal("Cannot stream results: " + err)
			// log.Printf("Today's date is: %d/%d/%d", month, day, year)
		}
		if err := resp.Error; err != nil {
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 {
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			}
			log.Fatalf("Could not recognize: %v", err)
		}
		for _, result := range resp.Results {
			fmt.Printf("Result: %+v\n", result)
		}
	}
}

// [END speech_transcribe_streaming_mic]

func textspeech() {
        // Instantiates a client.
        ctx := context.Background()

        client, err := texttospeech.NewClient(ctx)
        if err != nil {
                log.Fatal(err)
        }

        // Perform the text-to-speech request on the text input with the selected
        // voice parameters and audio file type.
        req := texttospeechpb.SynthesizeSpeechRequest{
                // Set the text input to be synthesized.
                Input: &texttospeechpb.SynthesisInput{
                        InputSource: &texttospeechpb.SynthesisInput_Text{Text: "Life is like a box of chocolates. You never know what you're gonna get."},
                },
                // Build the voice request, select the language code ("en-US") and the SSML
                // voice gender ("neutral").
                Voice: &texttospeechpb.VoiceSelectionParams{
                        LanguageCode: "en-US",
                        SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
                },
                // Select the type of audio file you want returned.
                AudioConfig: &texttospeechpb.AudioConfig{
                        AudioEncoding: texttospeechpb.AudioEncoding_MP3,
                },
        }

        resp, err := client.SynthesizeSpeech(ctx, &req)
        if err != nil {
                log.Fatal(err)
        }

        // The resp's AudioContent is binary.
        filename := "output.mp3"
        err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("Audio content written to file: %v\n", filename)
}