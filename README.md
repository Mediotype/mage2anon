# mage2anon

An easy-to-use tool to simultaneously dump and anonymize data for Magento 2. This would also work for Magento 1 and other applications as well.

## Usage

This tool uses a temporary file for the dump before anonymizing the data and removing the temporary file once complete.

```$xslt
mage2anon

    -config
        Defaults to the included configuration, otherwise should be the path for the configuration file you wish to use.
        
    -mysql-host
        MySQL Host - defaults to 127.0.0.1
        
    -mysql-user
        MySQL User - defaults to root
        
    -mysql-pass
        MySQL Password - defaults to root
        
    -mysql-port
        MySQL Port - defaults to 3306
        
    -mysql-db
        MySQL Database - *Required*
        
    -mysql-tables
        MySQL tables - defaults to nil, useful for small exports
```

## License
Copyright (c) 2019 [Mediotype](https://mediotype.com/).

[OSL 3.0](LICENSE.txt)

## Acknowledgements

This tool is based on the work done by Max Chadwick and his `dbanon` tool found [here](https://github.com/mpchadwick/dbanon).