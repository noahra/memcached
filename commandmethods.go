package main

import (
	"fmt"
)

func (cmd AppendCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd PrependCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd SetCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd GetCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd ReplaceCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd AddCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd PrependCommand) Execute(memcache *Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := parseDataBlock(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	currentValue, _ := memcache.Get(value.Key)
	value.DataBlock = append(dataBlock, currentValue.DataBlock...)
	amountOfBytes := value.AmountOfBytes
	value.AmountOfBytes = amountOfBytes + currentValue.AmountOfBytes

	_, ok := memcache.Get(cmd.key)
	if !ok {
		err := writeNotStored(cmd.connection)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	} else {
		memcache.Set(value.Key, value)
		err := noReplyCheck(cmd)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}

func (cmd AppendCommand) Execute(memcache *Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := parseDataBlock(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	currentValue, _ := memcache.Get(value.Key)
	value.DataBlock = append(currentValue.DataBlock, dataBlock...)
	amountOfBytes := value.AmountOfBytes
	value.AmountOfBytes = amountOfBytes + currentValue.AmountOfBytes
	_, ok := memcache.Get(cmd.key)
	if !ok {
		err := writeNotStored(cmd.connection)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	} else {
		memcache.Set(value.Key, value)
		err := noReplyCheck(cmd)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}

func (cmd SetCommand) Execute(memcache *Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := parseDataBlock(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	memcache.Set(value.Key, value)

	err = noReplyCheck(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

func (cmd AddCommand) Execute(memcache *Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := parseDataBlock(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	_, ok := memcache.Get(cmd.key)
	if ok {
		err := writeNotStored(cmd.connection)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	} else {

		memcache.Set(value.Key, value)
		err := noReplyCheck(cmd)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}

	return nil
}

func (cmd ReplaceCommand) Execute(memcache *Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := parseDataBlock(cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	_, ok := memcache.Get(cmd.key)

	if !ok {
		err := writeNotStored(cmd.connection)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	} else {
		memcache.Set(value.Key, value)
		err := noReplyCheck(cmd)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}

func (cmd GetCommand) Execute(memcache *Cache) error {
	val, ok := memcache.Get(cmd.key)
	if ok {
		if !ExpiryCheck(val) {
			//nolint:errcheck
			_, err := cmd.connection.Write([]byte(fmt.Sprintf("VALUE %s %d %d\n", val.Key, val.Flags, val.AmountOfBytes)))
			if err != nil {
				return fmt.Errorf("err: %w", err)
			}

			//nolint
			_, err = cmd.connection.Write([]byte(fmt.Sprintf("%s\n", val.DataBlock)))
			if err != nil {
				return fmt.Errorf("err: %w", err)
			}
		} else {
			memcache.Delete(cmd.key)
		}

		//nolint:errcheck
		_, err := cmd.connection.Write([]byte("END\r\n"))
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
	}
	return nil
}

func noReplyCheck(cmd Command) error {
	if cmd.GetBaseCommand().noReply != "noreply" {
		_, err := cmd.GetBaseCommand().connection.Write([]byte("STORED\r\n"))
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}
