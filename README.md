# zkPass-Node 
zkPass-Node is the server side implementation of zkPass protocol.
It performs multi-party secure computing, and can obtain protocol revenue by contributing verification computing power.

## Main Technologies

### **Multi-party computation**
Give users the ability to prove to third parties the provenance of such data. The specific implementation of this project, is based on:
- ECtF: Converting shares in EC(Fp) to shares in Fp
- GC: FreeXor + Half-Gates  (Stacked Garbling/Three Halves make a whole)
- OT: IKNP03 + KOS15 (Silent OT)

### **Zero-knowledge proofs**
Refer to the ability of a prover to convince a verifier that an assertion is correct without providing any useful information to the verifier.
The specific implementation of this project is based on the PLONK algorithm.

## Project Structure
The directory structure of the project looks like this:
```
├── app                  
│   ├── main.go          <- Main Entry
│   ├── cmd              <- All required circuits for GC
│   ├── ectf             <- ECtF implementation
│   ├── keystore         <- Manage the connection keys
│   ├── typings          <- Common Typings
│   ├── gc               <- Efficient garbling from a fixed-Key block cipher
│   │   ├── circuit      <- All required circuits for GC
│   │   ├── evaluator    <- GC evaluator P&P + FreeXor + Half-Gate
│   │   └── garbler      <- GC garbler
│   └── utils            <- Utilities
│
├── pkg                 
│   ├── chain            <- Read data and send transaction for chain 
│   ├── api              <- zkpass-node api
│   ├── node             <- zkpass-node main logic
│   ├── session          <- http session
│   ├── zkp              <- Verification of Zero knowledge Proof
│   └── ot               <- KOS15 Implementation
│
├── .gitignore            
├── go.mod               
├── LICENSE
└── README.md
```

## Quickstart
```bash
# make sure go env is installed on the system.

# clone project
`git clone https://github.com/zkPassOfficial/zkPass-node.git`
`cd zkPass-node`

# install dependencies
`go mod tidy`

# start with dev mode
`go run app/main.go start`

```

## Contributions
Have a question? Found a bug? Missing a specific feature? Feel free to file a new issue, discussion or PR with respective title and description.

Before making an issue, please verify that:

- The problem still exists on the current `main` branch.
- Your go dependencies are updated to recent versions.

Suggestions for improvements are always welcome!

## License
The zkPass-node library (i.e. all code in this repository) is licensed under the GNU GENERAL PUBLIC LICENSE Version 3, also included in our repository in the LICENSE file.