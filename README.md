# 🪐 Backend Engineer Challenge: Rockets 🚀

## Introduction 👋
Thank you for taking Lunar's code challenge for backend engineers! 

In the ZIP-file you have received, you will find a `README.md` (dah! of course) and folders 
containing executables for various operating systems and architectures.

> **Important:** If you cannot find an executable that works for you please reach out to us as soon as possible, 
> so we can get you one that works.

We hope you will enjoy this challenge - good luck.

## The Challenge 🧑‍💻
In this challenge you are going to build a service (or multiple) which consumes messages 
from a number of entities – i.e. a set of _rockets_ – and make the state of these 
available through a REST API. We imagine this API to be used by something like a dashboard.

As a minimum we expect endpoints which can:
1. Return the current state of a given rocket (type, speed, mission, etc.)
1. Return a list of all the rockets in the system; preferably with some kind of sorting.

The service should also expose an endpoint where the test program can post the messages to (see this [section](#running-the-test-program))

We are writing all our services in [Go](https://go.dev/) but there are no constrains on the programming language that you choose for 
solving the challenge. 
We prefer that you implement a great solution in a language that you feel comfortable in rather than trying to write 
in Go and implement a mediocre solution.

### The messages ✉️
Each rocket will be dispatching various messages (encoded as JSON) about its state changes through individual radio _channels_.
The channel is unique for each rocket and can therefore be treated as the ID of the rocket.

Apart from the channel each message also contains a _message number_ which expresses the order of the message within a channel, 
a _message time_ indicating when the message was sent and a _message type_ describing the event that occurred.

**Important:** Messages will arrive **out of order** and there is an **at-least-once guarantee** on messages 
meaning that you might receive the same message more than once.

Here is an example of a `RocketLaunch` message:

```json
{
    "metadata": {
        "channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
        "messageNumber": 1,    
        "messageTime": "2022-02-02T19:39:05.86337+01:00",                                          
        "messageType": "RocketLaunched"                             
    },
    "message": {                                                    
        "type": "Falcon-9",
        "launchSpeed": 500,
        "mission": "ARTEMIS"  
    }
}
```

The possible message types are:

#### `RocketLaunched`
Sent out once: when a rocket is launched for the first time.
```json
{
    "type": "Falcon-9",
    "launchSpeed": 500,
    "mission": "ARTEMIS"  
}
```

#### `RocketSpeedIncreased`
Continuously sent out: when the speed of a rocket is increased by a certain amount.
```json
{
    "by": 3000
}
```

#### `RocketSpeedDecreased`
Continuously sent out: when the speed of a rocket is decreased by a certain amount.
```json
{
    "by": 2500
}
```

#### `RocketExploded`
Sent out once: if a rocket explodes due to an accident/malfunction.
```json
{
    "reason": "PRESSURE_VESSEL_FAILURE"
}
```

#### `RocketMissionChanged`
Continuously sent out: when the mission for a rocket is changed.
```json
{
    "newMission":"SHUTTLE_MIR"
}
```

### Running the test program 💽
In the ZIP-file locate the executable that works for your system and run the following:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

This launches the program which starts posting (request method: `POST`) messages to the URL provided with a delay of 500ms between each message.

To see all commands run `./rockets help` and for help on the `launch` command run `./rockets launch --help`.

> We are going to run the program against your solution with the default values.

### Your solution and our assessment 📝
Before submitting your solution please make sure that you have included all the necessary files/information for 
running and assessing your solution. You can either submit a ZIP-file or provide a link to an online version control provider like GitHub, GitLab or Bitbucket.

Any design of a software system as a solution to a given problem will be affected by the choices made between various 
trade-offs. When submitting your solution, we will be really excited if you have explicitly described the design choices
you made and which trade-offs they entail. If you consciously chose a certain design, but are well aware that a different
(but maybe more complex) solution exist not having a certain trade-off, then include this in your documentation. 

When reviewing your solution we are going to look at things such as:
- The documentation provided, i.e. is it clear how to run your service(s) and, perhaps, what considerations/shortcuts have you made.
- The overall design of your solution, e.g. how easy is the code to understand, can the service(s) scale and how maintainable your code is.
- The measures you have taken to verify that your code works, e.g. automated tests.

We do not expect you to spend more than **6 hours** on this challenge. 
If you do not succeed in completing everything, then submit what you have, so we have something to look at - that is much better than nothing! ☺️
