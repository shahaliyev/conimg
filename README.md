# Concurrent Image Pixelizer

The program reads a .jpg file from path, after which, from left to right, top to bottom, and finds the average color for the (square size) x (square size) boxes. Then it sets the color of the whole square to that average color ([example](https://imgur.com/cvK7IxL)).

Two processing modes are possible: 
   - Single-processing [S]
   - Multi-processing [M]

The progress of pixelization is shown in a seperate window. After the complete rendering of the image, the result is stored in the ```result.jpg``` file. 

### The language of choice

For implementation, I chose Go programming language, as I had heard that it had a powerful support for concurrency. I did not get anything, but the app works fine.

### Installation

```$ go get github.com/shahaliyev/conimg```

### Application
The application will take three arguments from the command line: **file name, square size**, and the **processing mode**.

The program supports only **.jpg** file format and will terminate if the square size in pixels is negative or going out of the image's bounds. Two processing modes (single and multi-threaded) require the declaration of a case-sensitive single character as its argument -  **S** or **M**. Any other letter will terminate the program.

After go getting, you can run the following command line argument by specifying correct path for the image:

```$ go run main.go somefile.jpg 5 S```

### Medium

Read my medium article for explanation

###License

MIT