# go-fake-ts

A Fake time series data generator library written in Go (because it was an
opportunity to learn Go). 

This is the underlying library [TGEN](https://github.com/powerpu/tgen) uses to
generate the fake numbers.


## Background

I worked on a project where we needed to design a large time series data
warehouse whose input systems were unrelieable and many edge cases identified.
I even put together an article on the cases if you're interested
[here](https://www.linkedin.com/pulse/what-gotchas-edge-cases-when-processing-raw-time-series-dragan-rajak/).

The requirements for a fake timeseries generator were:

 * Generate data as a realtime stream or a batch input file

 * Generate plain text CSV or JSON in any format

 * Generate gigabytes of data for stress and volume testing

 * Generate millions of unique identities (see
   [here](https://www.linkedin.com/pulse/what-time-series-data-do-we-really-collect-dragan-rajak/)
   for background on terminology)

 * Generate consistent, repeatable data to facilitate testing

 * Simulate missing samples and missing data (two different concepts outlined
   below)

 * Generate seasonality, spikes, outliers and errors

I couldn't find what I wanted so as anyone would do I just put together
something quick and took it as an opportunity to learn a bit of Go.

This library only deals with generating the numbers in case you're feeling
frisky and want to use it in your own code. For a tool that implements all of
the above requirements use [TGEN](https://github.com/powerpu/tgen)


## Concepts


