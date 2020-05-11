package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Check if a record exists for locked FUND given an account address
func (k Keeper) AddressIsWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	if address.Empty() {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.WhitelistAddressStoreKey(address)
	return store.Has(addressKeyBz)
}

// AddAddressToWhitelist adds an address to the whitelist
func (k Keeper) AddAddressToWhitelist(ctx sdk.Context, address sdk.AccAddress) error {

	if address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.WhitelistAddressStoreKey(address), k.cdc.MustMarshalBinaryLengthPrefixed(address))

	return nil
}

// RemoveAddressFromWhitelist removes an address from the whitelist
func (k Keeper) RemoveAddressFromWhitelist(ctx sdk.Context, address sdk.AccAddress) error {

	if address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}

	if k.AddressIsWhitelisted(ctx, address) {
		store := ctx.KVStore(k.storeKey)
		store.Delete(types.WhitelistAddressStoreKey(address))
	}

	return nil
}

// IterateWhitelist iterates over the all the whitelisted addresses and performs a callback function
func (k Keeper) IterateWhitelist(ctx sdk.Context, cb func(addr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WhitelistKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addr sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &addr)

		if cb(addr) {
			break
		}
	}
}

// GetAllWhitelistedAddresses returns an array of all currently whitelisted addresses
func (k Keeper) GetAllWhitelistedAddresses(ctx sdk.Context) (addresses types.WhitelistAddresses) {
	k.IterateWhitelist(ctx, func(addr sdk.AccAddress) bool {
		addresses = append(addresses, addr)
		return false
	})
	return
}

// ProcessWhitelistAction processes the add/remove whitelist messages
func (k Keeper) ProcessWhitelistAction(ctx sdk.Context, address sdk.AccAddress, action types.WhitelistAction, signer sdk.AccAddress) error {

	logger := k.Logger(ctx)

	if address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}

	if !types.ValidWhitelistAction(action) {
		return sdkerrors.Wrap(types.ErrInvalidWhitelistAction, "action should be add or remove")
	}

	if !k.IsAuthorisedToDecide(ctx, signer) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist")
	}

	if action == types.WhitelistActionAdd {
		if !k.AddressIsWhitelisted(ctx, address) {
			err := k.AddAddressToWhitelist(ctx, address)
			if err != nil {
				return err
			}
			logger.Info("added address to purchase order whitelist", "address", address, "signer", signer)
		} else {
			return sdkerrors.Wrap(types.ErrAlreadyWhitelisted, fmt.Sprintf("%s already whitelisted", address))
		}
	}
	if action == types.WhitelistActionRemove {
		if k.AddressIsWhitelisted(ctx, address) {
			err := k.RemoveAddressFromWhitelist(ctx, address)
			if err != nil {
				return err
			}
			logger.Info("removed address from purchase order whitelist", "address", address, "signer", signer)
		} else {
			return sdkerrors.Wrap(types.ErrAddressNotWhitelisted, fmt.Sprintf("%s not whitelisted", address))
		}
	}

	return nil
}
