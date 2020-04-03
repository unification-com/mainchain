package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Check if a record exists for locked UND given an account address
func (k Keeper) AddressIsWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.WhitelistAddressStoreKey(address)
	return store.Has(addressKeyBz)
}

// AddAddressToWhitelist adds an address to the whitelist
func (k Keeper) AddAddressToWhitelist(ctx sdk.Context, address sdk.AccAddress) error {

	store := ctx.KVStore(k.storeKey)
	store.Set(types.WhitelistAddressStoreKey(address), k.cdc.MustMarshalBinaryLengthPrefixed(address))

	return nil
}

// RemoveAddressFromWhitelist removes an address from the whitelist
func (k Keeper) RemoveAddressFromWhitelist(ctx sdk.Context, address sdk.AccAddress) error {

	if k.AddressIsWhitelisted(ctx, address) {
		store := ctx.KVStore(k.storeKey)
		store.Delete(types.WhitelistAddressStoreKey(address))
	}

	return nil
}

// GetAllWhitelistedAddressesIterator returns an iterator for all current whitelisted addresses
func (k Keeper) GetAllWhitelistedAddressesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WhitelistKeyPrefix)
}

// GetAllWhitelistedAddresses returns an array of all currently whitelisted addresses
func (k Keeper) GetAllWhitelistedAddresses(ctx sdk.Context) types.WhitelistAddresses {
	whitelistIterator := k.GetAllWhitelistedAddressesIterator(ctx)

	var addresses types.WhitelistAddresses
	for ; whitelistIterator.Valid(); whitelistIterator.Next() {
		var addr sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(whitelistIterator.Value(), &addr)
		addresses = append(addresses, addr)
	}

	return addresses
}

func (k Keeper) ProcessWhitelistAction(ctx sdk.Context, address sdk.AccAddress, action types.WhitelistAction, signer sdk.AccAddress) error {

	logger := k.Logger(ctx)

	if !types.ValidWhitelistAction(action) {
		return sdkerrors.Wrap(types.ErrInvalidWhitelistAction, "action should be add or remove")
	}

	if !k.IsAuthorisedToDecide(ctx, signer) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist")
	}

	if action == types.WhitelistActionAdd {
		if !k.AddressIsWhitelisted(ctx, address) {
			_ = k.AddAddressToWhitelist(ctx, address)
			logger.Info("added address to purchase order whitelist", "address", address, "signer", signer)
		} else {
			return sdkerrors.Wrap(types.ErrAlreadyWhitelisted, fmt.Sprintf("%s already whitelisted", address))
		}
	}
	if action == types.WhitelistActionRemove {
		if k.AddressIsWhitelisted(ctx, address) {
			_ = k.RemoveAddressFromWhitelist(ctx, address)
			logger.Info("removed address from purchase order whitelist", "address", address, "signer", signer)
		} else {
			return sdkerrors.Wrap(types.ErrAddressNotWhitelisted, fmt.Sprintf("%s not whitelisted", address))
		}
	}

	return nil
}
