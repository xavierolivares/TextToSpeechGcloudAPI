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

import _ "github.com/joho/godotenv/autoload"

// [START speech_transcribe_streaming_mic]
import (
	"context"
	"fmt"
	"log"
	"io/ioutil"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

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
                        InputSource: &texttospeechpb.SynthesisInput_Text{Text: "Friends, I've got something new to tell all of you. I've decided to sponsor a hockey team made up entirely of chimps. I'm tired of people telling me that chimps are not capable of kicking human ass in sports. Chimps are just superior athletes. I've got them on strict diet of bone broth and elk meat."},
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