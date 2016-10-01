# Nojo

### Gateway for converting Noms to json writtern in Go

[Noms](https://github.com/attic-labs/noms) is a database with its own bespoke interface. Nojo lets you output Noms objects in json. It's set up to be deployed as a an AWS Lambda using [Gordon](https://github.com/jorgebastida/gordon).

### Build and run locally

Spit out a json object given a Noms path:
```js
$ go run nojo.go http://demo.noms.io/hn::raw.value.items[12211754]
{"by":"ahl","descendants":167,"id":12211754,"kids":[12214020,12213448,12212967,12212106,12212696,12212654,12212440,12212137,12211882,12213118,12212222,12215578,12212806,12212640,12212941,12212723,12215511,12215530,12214027,12213667,12215852,12212846,12212775,12216211,12231523,12212456,12215425,12212658,12213580,12212268,12213550,12219556,12212644,12222485,12216133,12211945,12212427,12216484],"score":508,"time":1470160624,"title":"Show HN: Noms – A new decentralized database based on ideas from Git","type":"story","url":"https://medium.com/@aboodman/noms-init-98b7f0c3566#.ojb6eaz94"}
```

We can use `json.tool` to help pretty this up:
```js
$ go run nojo.go http://demo.noms.io/hn::raw.value.items[10000000] | python -m json.tool
{
    "by": "rsp1984",
    "descendants": 1,
    "id": 10000000,
    "kids": [
        10003092
    ],
    "score": 37,
    "text": "",
    "time": 1438637523,
    "title": "Congratulations HN: 10M posts and comments",
    "type": "story",
    "url": "https://news.ycombinator.com/item?id=10000000"
}
```

### Build and deploy as an AWS Lambda

Nojo is designed to be used as an AWS Lambda. You'll need to [install Gordon](http://gordon.readthedocs.io/en/latest/installation.html). You'll also need to [install the AWS CLI tools](http://docs.aws.amazon.com/cli/latest/userguide/installing.html) and configure them.

First you'll need to specify the S3 bucket where you want your code to live.  By default, Gordon defines the bucket in the project's `settings.yml` file, but because the S3 bucket is global (per-region) sharing `settings.yml` files becomes tricky. Instead we specify our S3 bucket in an ancillary file and give it a (probably) unique name:

```sh
$ echo gordon-nojo-$(date +%s) >gordon/code-bucket
```

Now you can build and deploy the Lambda:

```sh
$ cd gordon
$ gordon build
Loading project resources
  ✓ apigateway:noms
Loading installed applications
  contrib_lambdas:
    ✓ lambdas:version
  nojo:
    ✓ lambdas:get
Building project...
  0001_p.json
  0002_pr_r.json
  0003_r.json
qiviut /Users/ahl/src/nojo/gordon $ gordon apply
Applying project...
  0001_p.json (cloudformation)
    ✓ No updates are to be performed.
  0002_pr_r.json (custom)
    ✓ code/contrib_lambdas_version.zip (9fb42e49)
    ✓ code/nojo_get.zip (d040ff84)
  0003_r.json (cloudformation)
    UPDATE_COMPLETE_CLEANUP_IN_PROGRESS waiting... |
Project Outputs:
  LambdaNojoGet
    arn:aws:lambda:us-west-2:935244743057:function:dev-nojo-r-NojoGet-188JR8YZYZ903:current
  ApigatewayNoms
    https://9tvp9qurzi.execute-api.us-west-2.amazonaws.com/dev
```

Then you can use your browser or curl (or your application!) to access the data:

```js
$ curl -s "https://9tvp9qurzi.execute-api.us-west-2.amazonaws.com/dev?nomspath=http://demo.noms.io/hn::raw.value.items\[3675309\]" | python -m json.tool
{
    "by": "moadeel",
    "descendants": 1,
    "id": 3675309,
    "kids": [
        3675389
    ],
    "score": 2,
    "text": "I have seen a lot of pictures trying to depict the \"mightiness\" of humans against smaller organisms or their \"tinyness\" compared to celestial bodies ... but for me this one takes the cake - be sure to give it a shot.",
    "time": 1331125819,
    "title": "Awesome: how humans compare with other beings in size",
    "type": "story",
    "url": "http://htwins.net/scale2/scale2.swf?bordercolor=white&fb_source=message"
}
```

#### Contributing

Contributions of all kinds welcome: issues, questions, docs, code.
