import React from "react";
import { Table } from "reactstrap";

const DomainTable = (props) => {
  function GetRows() {
    return props.rows.map((item, key) => {
      return (
        <tr key={key}>
          <td>{item.domain}</td>
          <td>{item.validFrom}</td>
        </tr>
      );
    });
  }

  return (
    <Table striped>
      <thead>
        <tr>
          <th>Domain</th>
          <th>Valid From</th>
        </tr>
      </thead>
      <tbody>
        <GetRows></GetRows>
      </tbody>
    </Table>
  );
};

export default DomainTable;
