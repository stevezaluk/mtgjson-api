<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/stevezaluk/mtgjson-api">
    <img src="docs/images/logo-mtgjson.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">MTGJSON-API</h3>

  <p align="center">
    A RESTful API for interacting with Magic: The Gathering card data
    <br />
    <a href="https://github.com/stevezaluk/mtgjson-api"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://stevezaluk.atlassian.net/jira/software/projects/SCRUM/boards/1/backlog">Jira Board</a>
    ·
    <a href="https://stevezaluk.atlassian.net/jira/software/projects/SCRUM/boards/1/backlog">Report Bug</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li><a href="#disclaimers">Disclaimers</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#configuration">Configuration</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

MTGJSON-API is a RESTful API written in Go, built ontop of the MTGJSON dataset. It features full integration with Auth0 to provide authentication and will eventually support fetching prices as well. Additionally, it allows users to create there own deck and fetch them through the API. This was originally built for the in-progress MTG Simulator: Arcane, however it was developed separately so that you can run this without running Arcane.

## Disclaimers

MTGJSON API is unofficial Fan Content permitted under the Fan Content Policy. Not approved/endorsed by Wizards of the Coast. Portions of the materials used are property of Wizards of the Coast. © Wizards of the Coast LLC.

MTGJSON API not officially endorsed by MTGJSON

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Getting started with the API is fairly simple, provided that the following pre-requisites below are deployed and running. MTGJSON-API relies on three technologies to provide its functionality:
* Go Version 1.23.2 
* MongoDB
* Auth0 (optional)

You can view a very basic guide below to get started with these, however you should plan your own deployment of these independently and should not rely on the below guides for starting a production environment

### Prerequisites
To start, install Go Version 1.23.2 using the guide present on Golang's website: https://go.dev/doc/install. After this is done MongoDB should be installed and configured. In this tutorial we are going to use docker to deploy this quickly however, this does not represent a properly deployed production environment. Additionally, instructions on how to install docker is out of the scope of this tutorial, and assumes that docker is fully installed and functional

* MongoDB
  ```sh
  docker pull mongo@latest
  docker run --name some-mongo -d mongo:tag -p 27017:27017
  ```

Auth0 is used as an authentication provider for the API and is optional to use. You can create a free account with 1000 tokens per month alloted to you at auth0.com. For more specifics on how to further configure this for use with MTGJSON-API, please see the guide below

### Configuration

You can define the mandatory configuration values through three different methods: A JSON config file, environmental variables, or through the CLI flags present in the application. Cobra is used as the CLI for the API, so any configuration values defined in a config file can also be defined using the proper flags. The API checks the following path for the config file by default: ~/.config/mtgjson-api/config.json

#### Mandatory Flags

The following configuration values are mandatory for starting the API server:

* Mongo DB IP Address (string) ```mongo.ip``` - The IP Address of your running MongoDB instance
* Mongo DB Port (integer) ```mongo.port``` - The port your MongoDB instance is listening for connections on
* Mongo DB Username (string)  ```mongo.user``` - The username to use for authentication with MongoDB
* Mongo DB Password (string) ```mongo.pass``` - The password to use for authentication with MongoDB

#### Authentication Flags
If you wish to enable authentication within the API the following values must also be set:

* Auth0 Domain (string) ```auth0.domain``` - The domain of your Auth0 tenant
* Auth0 Audience (string) ```auth0.audience``` - The identifier (audience) of your API
* Auth0 Client ID (string) ```auth0.client_id``` - The Client ID for your Auth0 Application
* Auth0 Client Secret (string) ```auth0.client_secret``` - The Client secret for your Auth0 application

Additionally, if you decide to run the API without authentication, you can bypass with the following values. These can also be defined as flags which can be useful for testing the API:

* API No Auth (bool) ```api.no_auth``` - Disable authentication for all endpoints
* API No Scope (bool) ```api.no_scope``` - Disable scope validation for all endpoints

#### Log Flags

Finally you can define the path in which log files are stored using the following flag:

* Log Path (string) ```log.path``` - The unix path to store log files in

### Auth0 Configuration

To properly configure authentication with Auth0, a few things need to be completed within your Auth0 tenant. Below is a step by step tutorial on how to complete these steps. Please be aware that screenshots are not included here, and if there is any confusion please consult Auth0's documentation

#### Application Setup

1. Create a free Auth0 account at https://www.auth0.com
2. Create a new Machine to Machine application in your Auth0 tenant. Make note of your Client ID and Client Secret
3. Click settings under your application, scroll all the way down and select the drop down for Advanced Settings
4. Click Grant Types, and select Client Credentials, Password, and MFA
5. Click Save Settings.

#### API Setup

1. Create a new API, name this however best suits your needs and create an identifier you will remember
2. Click Settings, and scroll down to "RBAC Settings"
3. Click "Enable RBAC", and then click "Add Permissions in Access Token"
4. Click Machine to Machine Applications and ensure that it says "Authorized"
5. Click Save Settings.

#### Permissions (Scope Setup)

A specific set of scopes must be defined for the API you have created for the scope validation to function properly. Below is a list of scoped permissions that must be created:

* Read set permissions ```read:set``` - Provides permissions for indexing all sets and reading metadata from individual ones
* Write set permissions ```write:set``` - Provides permissions for creating, modifying, and deleting sets
* Read card permissions ```read:card``` - Provides permissions for indexing all cards and reading metadata from individual cards
* Write card permissions ```write:card``` - Provides permissions for creating, modifying, and deleting cards
* Read deck permissions ```read:deck``` - Provides permissions for indexing all decks and reading metadata from individual ones
* Write deck permissions ```write:deck``` - Provides permissions for creating, modifying, and deleting decks
* Read user permissions ```read:user``` - Provides permissions for indexing all users and reading metadata from individual ones
* Write user permissions ```write:user``` - Provides permissions for creating new users and modifying existing ones
* Read health permissions ```read:health``` - Provides permissions for fetching the health status of the API
* Read metric permissions ```read:metrics``` - Provides permissions for fetching prometheus metrics from the API

To add these permissions to the Auth0 API, follow the steps below:

1. Click Applications and then API's on the left sidebar
2. Click the API you created in the previous steps
3. Click "Permissions"
4. In the "Permission" box, define the permission as listed below and attach a description for it
5. Click the "Add" button on the right of the "Description" box
6. Repeat steps 4-5 for the rest of the permissions

#### Final Steps

After the above steps have been completed, make note of your Client ID, Client Secret, Domain, and Audience (API Identifier) and ensure that you add them to your configuration file or pass them through CLI Flags

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/stevezaluk/mtgjson-api.git
   ```
2. Install dependencies
   ```sh
   go get .
   ```
3. Define your configuration values as described above
4. Build the project
   ```sh
    go build
   ```
5. Run the API
    ```sh
    ./mtgjson run
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under Apache License 2.0. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Steven A. Zaluk - [@steve_zaluk](https://x.com/stevezaluk) - arcanegame@protonmail.com

Project Link: [https://github.com/stevezaluk/mtgjson-api](https://github.com/stevezaluk/mtgjson-api)\n

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/stevezaluk/mtgjson-api.svg?style=for-the-badge
[contributors-url]: https://github.com/stevezaluk/mtgjson-api/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/stevezaluk/mtgjson-api.svg?style=for-the-badge
[forks-url]: https://github.com/stevezaluk/mtgjson-api/network/members
[stars-shield]: https://img.shields.io/github/stars/stevezaluk/mtgjson-api.svg?style=for-the-badge
[stars-url]: https://github.com/stevezaluk/mtgjson-api/stargazers
[issues-shield]: https://img.shields.io/github/issues/stevezaluk/mtgjson-api.svg?style=for-the-badge
[issues-url]: https://github.com/stevezaluk/mtgjson-api/issues
[license-shield]: https://img.shields.io/github/license/stevezaluk/mtgjson-api.svg?style=for-the-badge
[license-url]: https://github.com/stevezaluk/mtgjson-api/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/stevezaluk/
[go-sdk-version]: https://img.shields.io/github/go-mod/go-version/stevezaluk/mtgjson-sdk
