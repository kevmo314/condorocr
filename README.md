# kevmo314/condorocr

Performs OCR for Condor frames.

Get some frames:

```
ffmpeg -i source.mp4 -r 1 output_%04d.png
```

Then run the ocr

```
cat output_0001.png | go run main.go
```