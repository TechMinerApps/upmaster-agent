# Basic Requirements

## Basic Functions
- Poll endpoint list from UpMaster
- Check endpoint status
- Write data back to UpMaster

## Platform Requirements

Require Go >= 1.15.8

## Commit Message Convention

Use AngularJS style commit message.

# Development

## Rough Structure

Endpoint Poller <- UpMaster

Endpoint Pool

Endpoint Checker -> Data Writer -> UpMaster

## Detailed Design

Communication between components is done by passing message through channel.

### Endpoint Pool

Endpoint Pool contains all the endpoints and information required for other components. A Mutex Lock is required for the pool.

### Endpoint Poller

Endpoint Poller get the endpoint list from UpMaster and refresh the Endpoint Pool every specific time.

### Endpoint Checker

Each endpoint has one checker to deal with check interval. After checking, a message containing information is pass to Data Writer by channel.

### Data Writer

Data Writer accepts data from a queue and perform write operation through UpMaster API