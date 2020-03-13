# Autonity Documentation Style Guide

## Purpose of this Document

Follow these guidelines to ensure the Autonity documentation is consistent and well organised.

This is a living document and will evolve to better suit Autonity users and contributors needs.

> **Note:** Although not everything in this style guide is currently followed in the Autonity
documentation, these guidelines are intended to be applied when writing new content and revising
existing content.

**The primary audience for this document is:**

*   Members of the Autonity team
*   Developers and technical writers contributing to the Autonity documentation

## Mission Statement

The Autonity documentation contributes to a consistent and a great experience for end users of Ethereum clients.

## General Guidelines

The guiding principles for the Autonity documentation are:
1. Consistent
1. Simplity and technical accuracy
1. Proactivity and good practice
1. Informative and exhaustive

### 1. Consistency

Consistency is important to help our end users build a conceptual model of how Autonity works.
Consistency with words, formatting, and style helps users know what to expect when they refer to or search Autonity documentation.

### 2. Simpicity and technical accuracy

Avoid jargon and always assume our end users may not be Ethereum experts.

Explain Autonity functionality and when an understanding of complex Ethereum concepts is required refer users to relevant resources.

For example, to explain how the EVM works, link to ethdocs.org documentation such as
https://github.com/ethereum/wiki/wiki/Ethereum-Virtual-Machine-(EVM)-Awesome-List

Simple explanations must still be technically correct.

### 3. Proactivity and good practices

Being proactive means anticipating user needs and guiding them through a process.
This most often takes the form of notes or tip messages alongside the main explanation.
Put yourself in the place of a user and consider what questions you would have when reading the documentation.

Do not assume required steps are implied. Include them if you are unsure.

Documenting good practices is also important.
For example, instruct users to secure private keys and protect RPC endpoints a production environments.

### 4. Informative and exhaustive

We seek a balance between providing enough relevant information to help our users develop a solid
conceptual model of how Autonity works without forcing them to read too much text or redundant detail.

To provide additional detail, use sub-sections.

## Writing Style Guide

Here are some important points we follow:

### Active Voice
Use active voice. Use _you_ to create a more personal friendly style. Avoid gendered pronouns (_he_ or _she_).

### Contractions
Use contractions. For example, don’t.

Use common contractions, such as it’s and you’re, to create a friendly, informal tone.

### Recommend
It's acceptable to use "we recommend" to introduce a product recommendation.
Don't use "Autonity recommends" or "it is recommended."

Example: Instead of _This is not recommended for production code_ use _We don't recommend this for production code_.

### Directory vs Folder
Use _directory_ over _folder_ because we are writing for developers.

### Headings
Captitalise the first word of headings.


### Assumed Knowledge For Readers
We have two distinct audiences to consider when developing content:

- New to Ethereum and Ethereum clients
- Experienced with Ethereum clients other than Autonity.

### Avoid Abbreviations

Try not to use abbreviations [except for well known ones and some jargon](MKDOCS-MARKDOWN-GUIDE.md#abbreviations).
Don't use "e.g." but use "for example".
Don't use "i.e." but use "that is".
