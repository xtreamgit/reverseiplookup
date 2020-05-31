import React, { useState } from "react";
import { Container, Row, Col } from "reactstrap";
import SearchBar from "./components/SearchBar.js";
import DomainTable from "./components/DomainTable.js";
import ApiHandler from "./apiHandler.js";
require("dotenv").config();

function App() {
  const [domainRows, setDomainRows] = useState([]);

  const handleSearchClick = (e) => {
    //ip regex also validates for host.
    if (ApiHandler.validateIP(e.searchTerm)) {
      ApiHandler.searchByIP(e.searchTerm).then((res) => {
        console.log(res);
        setDomainRows(res.message);
      });
    } else if (ApiHandler.validateHost(e.searchTerm)) {
      ApiHandler.searchByDomain(e.searchTerm).then((res) => {
        console.log(res);
        setDomainRows(res.message);
      });
    } else {
      console.log(e, "couldnt regex string");
    }
  };

  return (
    <Container>
      <Row>
        <Col>
          <h2>Reverse IP Lookup</h2>
        </Col>
      </Row>
      <Row>
        <Col>
          <SearchBar onClick={handleSearchClick}></SearchBar>
          <DomainTable rows={domainRows}></DomainTable>
        </Col>
      </Row>
    </Container>
  );
}

export default App;
