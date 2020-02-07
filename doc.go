/*
Package fake allows for generation of fake times, true/false values and numbers
that should follow a predefined pattern. The intent is to generate time series
data.
*/

package fake

/*
Concepts

The package breaks down time series data into two major parts, samples and the
actual data.

There is a distinctinction between "bad samples" and "bad data". 

An example of the distinction would be in an example of collecting CPU and
memory data from a server.

In a "valid sample" the collection server would connect and retrieve this
information from the destination server in a particular way with no issues.

In a "bad sample" scenario, the destination server would be running fine but
due to a network outage a sample would not be collected.

In a "bad data" scenario, the destination server would be running and the
network would be up but due to a rogue process the CPU value would not be able
to be collected while memory could. The end result is that the CPU data is
"bad".

Basics

There are 4 types in this package and they all work in the same way. They are:

  type Pattern struct {}
  type Random struct {}
  type Time struct {}
  type Data struct {}

Each implements the following interface for convenience:

  type Value interface{}

The general usage would be to instantiate one or more types with

  fakeType := New<Pattern|Random|Time|Data>()

and then call fakeType.Next() to generate a next value once and
fakeType.Value() to retrieve it one or more times.

The goal is to generate time series data that may look as follows:

  Timestamp, CPU, Memory
  1581039550,23.5,97.2
  1581039650,26.5,84.9
  ...etc...

One would use the true/false value of a Pattern or Random to determine what to
print.

An example simulation scenario could be generating data where every 24th hourly
sample fails because you want to simulate a scheduled job taking up too many
resources and thus not retrieving data.

For this you would use a 23 good samples and 1 bad samples pattern with time
going up by 1 hour intervals. Your pseoudocode may look as follows:

  fakePattern := fake.NewPattern(...)
  fakeTime := fake.NewTime(...)
  fakeData1 := fake.NewData(...)
  fakeData2 := fake.NewData(...)

  for i := 0; i < 1000; i++ {
      if fakePattern.Value() {
          fmt.Printf("%v: %v,%v,%v\n", i, fakeTime.Value(), fakeData1.Value(), fakeData2.Value())
      } else {
          fmt.Printf("Bad sample at %v!\n", i, fakeTime.Value())
      }
  }

If you want to get fancy you could then simulate a bad network connection that
drops away randomly 5% of the time. You would modify the above as follows:

  fakePattern := fake.NewPattern(...)
  fakeRandom := fake.NewRandom(...)
  fakeTime := fake.NewTime(...)
  fakeData1 := fake.NewData(...)
  fakeData2 := fake.NewData(...)

  for i := 0; i < 1000; i++ {
      if fakePattern.Value() | fakeRandom.Value() {
          fmt.Printf("%v: %v,%v,%v\n", i, fakeTime.Value(), fakeData1.Value(), fakeData2.Value())
      } else {
          fmt.Printf("Bad sample at %v!\n", i, fakeTime.Value())
      }
  }

You can go as wild or as simple as you like. Remember, you're only limited by
your imagination!
*/
