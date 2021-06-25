## insprctl delete apps

Delete apps from scope

```
insprctl delete apps [flags]
```

### Examples

```
  # Delete app from the default scope
 insprctl delete apps <appname> 

  # Delete app from a custom scope
 insprctl delete apps <appname> --scope app1.app2

```

### Options

```
  -c, --config string   set the config file for the command
  -h, --help            help for apps
  -s, --scope string    insprctl <command> --scope app1.app2
  -t, --token string    set the token for the command
```

### SEE ALSO

* [insprctl delete](insprctl_delete.md)	 - Delete component of object type

###### Auto generated by spf13/cobra on 15-Jun-2021