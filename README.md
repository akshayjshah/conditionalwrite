# Conditional S3 writes with `aws-sdk-go-v2`

A small example repository showing how to issue [conditional writes][] with v2
of AWS's Go SDK and run integration tests with MinIO. There's a [short
discussion on my blog][blog].

**With the Docker daemon running** (or socket-activated), run the tests with
`make test` or `go test .`.

> [!IMPORTANT]
> **Don't import this code.** Copy it into your project and customize it instead!
> It's under the [MIT License](LICENSE), so you should be able to use it
> nearly anywhere.

[conditional writes]: https://aws.amazon.com/about-aws/whats-new/2024/08/amazon-s3-conditional-writes/
[blog]: https://akshayshah.org/s3-conditional-writes-go-sdk-v2
