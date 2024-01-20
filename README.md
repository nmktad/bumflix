bumflix. soon to be the best streaming service in the world

to run it locally. clone the project put any video you want to stream in the partition directory and run the command below in the same directory

```bash
ffmpeg -i <videoname> -profile:v baseline -level 3.0 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls index.m3u8
```

if you don't have ffmpeg installed run the command below

```bash
sudo apt install ffmpeg
```

to run the app you can use the command below in the root directory
```bash
go run .
```

clone the client app from [bumflixweb](https://github.com/nmktad/bumflixweb) and run it to see the stream in a nextjs web app
