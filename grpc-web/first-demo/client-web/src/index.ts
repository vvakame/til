import { ChatServiceClient } from '../generated/chat_grpc_web_pb';
import { SpeakRequest, SpeakResponse, ListenSpeakRequest } from '../generated/chat_pb';

const chatService = new ChatServiceClient('http://localhost:8080', {}, {});

const request = new SpeakRequest();
request.setUserName("vvakame");
request.setMessage('Hello World!');

{
    const request = new ListenSpeakRequest();
    const stream = chatService.listenSpeak(request, {})
    stream.on('data', response => {
        console.log(response);
    });
}
{
    const call = chatService.speak(request, {}, (err, resp) => {
        if (err) {
            console.error(err);
            return;
        }
        console.log(resp);
    });
    call.on('status', status => {
        console.log(status);
    });
}
