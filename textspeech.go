
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
                        InputSource: &texttospeechpb.SynthesisInput_Text{Text: "I see trees of green, red roses too. I see them bloom for me and you. And I think to myself what a wonderful world. I see skies of blue and clouds of white. The bright blessed day, the dark sacred night. And I think to myself what a wonderful world"},
                },
                // Build the voice request, select the language code ("en-US") and the SSML
                // voice gender ("neutral").
                Voice: &texttospeechpb.VoiceSelectionParams{
                        LanguageCode: "en-US",
                        SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
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