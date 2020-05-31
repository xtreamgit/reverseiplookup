import React, { useState } from "react";
import { InputGroup, InputGroupAddon, Input, Button } from "reactstrap";

const SearchBar = (props) => {
  const [values, setValues] = useState({ searchTerm: "" });
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setValues({ ...values, [name]: value });
  };
  const handleClick = (e) => {
      props.onClick(values)
  };

  return (
    <InputGroup>
      <Input
        name="searchTerm"
        onChange={handleInputChange}
        value={values.searchTerm}
      />
      <InputGroupAddon addonType="prepend">
        <Button onClick={handleClick}>Search</Button>
      </InputGroupAddon>
    </InputGroup>
  );
};

export default SearchBar;
