package main

import (
	"fmt"
	"strconv"
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
	value.DataBlock = dataBlock + currentValue.DataBlock
	amountOfBytes := value.AmountOfBytes
	value.AmountOfBytes = amountOfBytes + currentValue.AmountOfBytes

	_, ok := memcache.Get(cmd.key)
	if !ok {
		writeNotStored(cmd.connection)
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
	value.DataBlock = currentValue.DataBlock + dataBlock
	amountOfBytes := value.AmountOfBytes
	value.AmountOfBytes = amountOfBytes + currentValue.AmountOfBytes
	_, ok := memcache.Get(cmd.key)
	if !ok {
		writeNotStored(cmd.connection)
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
		writeNotStored(cmd.connection)
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
		writeNotStored(cmd.connection)
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
			cmd.connection.Write([]byte("VALUE " + val.Key + " " + val.Flags + " " + strconv.Itoa(val.AmountOfBytes) + "\n"))
			cmd.connection.Write([]byte(val.DataBlock + "\n"))
		} else {
			memcache.Delete(cmd.key)
		}
		cmd.connection.Write([]byte("END\r\n"))
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
