### Streaming vs. Downloading
- Downloading is when you wait for the entire file to be transferred before you can start using it.
- Streaming is when you start using the file immediately while it's still being transferred in the background.

### The simplest way to stream a video file on the web (imo) is to take advantage of two things:
- the native HTML5 <video> element. It streams video files by default as long as the server supports it.
- the Range HTTP header. It allows the client to request specific byte ranges of a file, enabling partial downloads. S3 servers support it by default.
- Writing streaming from scratch is hard, but using the right tools makes it pretty easy these days.

> NOTE: ```206``` is the HTTP status code for "Partial Content", which is what the server responds with when it successfully serves a byte range of a file.


- Well, in a "traditional" mp4 file (as our current Boots videos are), the "moov" atom (which contains metadata about the video) is at the end of the file.
- But the client needs that metadata before it can start playing the video (metadata contains details like video duration, frame rate, encoding format, decoding format, etc.). So the browser is smartly poking around using Range requests to get that metadata as quickly as it can. 
- We can speed up this process by pre-processing the video to have "fast start" encoding by moving the moov atom to the start.

### EXTRA INFO
- Fast start moov: If the moov atom is at the beginning of the file (before the actual video data), the browser can immediately read the metadata and start playback without waiting for the entire file to download.
- Standard moov: If the moov atom is at the end of the file, the browser must download much more of the file before it knows the video's properties and can begin playing.

### Why metadata is at the end of the file in standard mp4 files
- The moov atom is at the end of the file because it contains information about the entire video, including the duration and the locations of all the frames. When the video is being encoded, the encoder doesn't know this information until it has processed the entire video. So it writes the moov atom at the end of the file after all the video data has been written.

### How to make a video file "fast start"
- This is why fast-start files exist—tools like FFmpeg can reorganize files after encoding to move the moov atom to the front. However, this:
```
Adds an extra processing step (slower)
Requires reading and rewriting the entire file
Isn't always done automatically because it costs time and resources
```
- So it's a trade-off: standard encoding is faster and simpler, but fast-start optimization improves playback but requires extra work afterward.

### Extra Knowledge
> Adaptive streaming: Standard mp4 files have a single resolution and bitrate. If a user's connection speed is unstable, HLS or MPEG-DASH allows for changing the quality of the stream on the fly. You may have noticed on YouTube or Netflix that your video quality changes based on your connection speed. Dropping to lower resolution is better than endlessly buffering.

> Live streaming: Standard mp4 files are not designed to be updated in real-time. You'd want to use a lower-latency protocol like WebRTC or RTMP for live streaming.