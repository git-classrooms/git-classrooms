import React from "react";
import { Button } from "../../components/index";

export function Sidebar(): JSX.Element {
  return (
    <div>
      <Button
        onClick={() => console.log("clicked the test button")}
        text="Test-Button"
      />
    </div>
  );
}
