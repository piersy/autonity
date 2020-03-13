Generating private and public keys
-----------------------------------

Using OpenSSL to generate a private and public keys.

.. note:: Store these in a directory somewhere safe

1. Generate a private key.

	.. code-block:: bash

		openssl genrsa -out privatekey.pem 1024


2. Create a X509 certificate (.cer file) containing your public key.

	.. code-block:: bash

		openssl req -new -x509 -key privatekey.pem -out publickey.cer -days 1825



3. Export your x509 certificate and private key to a pfx file.

	.. code-block:: bash

		openssl pkcs12 -export -out public_privatekey.pfx -inkey privatekey.pem -in publickey.cer
