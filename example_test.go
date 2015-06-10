package xmlprocessor

import (
	"bytes"
	"encoding/xml"
	"os"
)

func Example_1Issue6800() {
	// https://github.com/golang/go/issues/6800#issuecomment-83049402
	example := `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>`

	src := bytes.NewBufferString(example)
	NewDefaultProcessor().ProcessStreams(os.Stdout, src)
	// Output:
	// <w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>
}

func Example_2Issue6800() {
	// https://github.com/golang/go/issues/6800#issuecomment-82772038
	example := `<taxii_11:Discovery_Response 
    xmlns:taxii="http://taxii.mitre.org/messages/taxii_xml_binding-1"
    xmlns:taxii_11="http://taxii.mitre.org/messages/taxii_xml_binding-1.1"
    xmlns:tdq="http://taxii.mitre.org/query/taxii_default_query-1" 
    message_id="32898" 
    in_response_to="1">
</taxii_11:Discovery_Response>`

	src := bytes.NewBufferString(example)
	dec := xml.NewDecoder(src)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	defer enc.Flush()

	NewDefaultProcessor().Process(enc, dec)
	// Output:
	// <taxii_11:Discovery_Response xmlns:taxii="http://taxii.mitre.org/messages/taxii_xml_binding-1" xmlns:taxii_11="http://taxii.mitre.org/messages/taxii_xml_binding-1.1" xmlns:tdq="http://taxii.mitre.org/query/taxii_default_query-1" message_id="32898" in_response_to="1"></taxii_11:Discovery_Response>
}

func Example_3Issue6800() {
	// https://github.com/golang/go/issues/6800#issuecomment-76655588
	example := `
<stix:STIX_Package
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:stix="http://stix.mitre.org/stix-1"
xmlns:indicator="http://stix.mitre.org/Indicator-2"
xmlns:cybox="http://cybox.mitre.org/cybox-2"
xmlns:AddressObject="http://cybox.mitre.org/objects#AddressObject-2"
xmlns:cyboxVocabs="http://cybox.mitre.org/default_vocabularies-2"
xmlns:stixVocabs="http://stix.mitre.org/default_vocabularies-1"
xmlns:ExampleNamespace="http://example.com/"
xsi:schemaLocation="
http://stix.mitre.org/stix-1 http://stix.mitre.org/XMLSchema/core/1.0.1/stix_core.xsd
http://stix.mitre.org/Indicator-2 http://stix.mitre.org/XMLSchema/indicator/2.0.1/indicator.xsd
http://cybox.mitre.org/default_vocabularies-2 http://cybox.mitre.org/XMLSchema/default_vocabularies/2.0.1/cybox_default_vocabularies.xsd
http://stix.mitre.org/default_vocabularies-1 http://stix.mitre.org/XMLSchema/default_vocabularies/1.0.1/stix_default_vocabularies.xsd
http://cybox.mitre.org/objects#AddressObject-2 http://cybox.mitre.org/XMLSchema/objects/Address/2.0.1/Address_Object.xsd"
id="ExampleNamespace:STIXPackage-33fe3b22-0201-47cf-85d0-97c02164528d"
version="1.0.1">
<stix:STIX_Header>
    <stix:Title>Example watchlist that contains IP information.</stix:Title>
    <stix:Package_Intent xsi:type="stixVocabs:PackageIntentVocab-1.0">Indicators - Watchlist</stix:Package_Intent>
</stix:STIX_Header>
<stix:Indicators>
    <stix:Indicator xsi:type="indicator:IndicatorType" id="ExampleNamespace:Indicator-33fe3b22-0201-47cf-85d0-97c02164528d">
        <indicator:Type xsi:type="stixVocabs:IndicatorTypeVocab-1.0">IP Watchlist</indicator:Type>
        <indicator:Description>Sample IP Address Indicator for this watchlist. This contains one indicator with a set of three IP addresses in the watchlist.</indicator:Description>
        <indicator:Observable  id="ExampleNamespace:Observable-1c798262-a4cd-434d-a958-884d6980c459">
            <cybox:Object id="ExampleNamespace:Object-1980ce43-8e03-490b-863a-ea404d12242e">
                <cybox:Properties xsi:type="AddressObject:AddressObjectType">
                    <AddressObject:Address_Value condition="Equals" apply_condition="ANY">10.0.0.0##comma##10.0.0.1##comma##10.0.0.2</AddressObject:Address_Value>
                </cybox:Properties>
            </cybox:Object>
        </indicator:Observable>
    </stix:Indicator>
</stix:Indicators>
</stix:STIX_Package>
`

	src := bytes.NewBufferString(example)
	dec := xml.NewDecoder(src)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	defer enc.Flush()

	NewDefaultProcessor().Process(enc, dec)
	// Output:
	// <stix:STIX_Package xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:stix="http://stix.mitre.org/stix-1" xmlns:indicator="http://stix.mitre.org/Indicator-2" xmlns:cybox="http://cybox.mitre.org/cybox-2" xmlns:AddressObject="http://cybox.mitre.org/objects#AddressObject-2" xmlns:cyboxVocabs="http://cybox.mitre.org/default_vocabularies-2" xmlns:stixVocabs="http://stix.mitre.org/default_vocabularies-1" xmlns:ExampleNamespace="http://example.com/" xsi:schemaLocation="&#xA;http://stix.mitre.org/stix-1 http://stix.mitre.org/XMLSchema/core/1.0.1/stix_core.xsd&#xA;http://stix.mitre.org/Indicator-2 http://stix.mitre.org/XMLSchema/indicator/2.0.1/indicator.xsd&#xA;http://cybox.mitre.org/default_vocabularies-2 http://cybox.mitre.org/XMLSchema/default_vocabularies/2.0.1/cybox_default_vocabularies.xsd&#xA;http://stix.mitre.org/default_vocabularies-1 http://stix.mitre.org/XMLSchema/default_vocabularies/1.0.1/stix_default_vocabularies.xsd&#xA;http://cybox.mitre.org/objects#AddressObject-2 http://cybox.mitre.org/XMLSchema/objects/Address/2.0.1/Address_Object.xsd" id="ExampleNamespace:STIXPackage-33fe3b22-0201-47cf-85d0-97c02164528d" version="1.0.1">
	//   <stix:STIX_Header>
	//     <stix:Title>Example watchlist that contains IP information.</stix:Title>
	//     <stix:Package_Intent xsi:type="stixVocabs:PackageIntentVocab-1.0">Indicators - Watchlist</stix:Package_Intent>
	//   </stix:STIX_Header>
	//   <stix:Indicators>
	//     <stix:Indicator xsi:type="indicator:IndicatorType" id="ExampleNamespace:Indicator-33fe3b22-0201-47cf-85d0-97c02164528d">
	//       <indicator:Type xsi:type="stixVocabs:IndicatorTypeVocab-1.0">IP Watchlist</indicator:Type>
	//       <indicator:Description>Sample IP Address Indicator for this watchlist. This contains one indicator with a set of three IP addresses in the watchlist.</indicator:Description>
	//       <indicator:Observable id="ExampleNamespace:Observable-1c798262-a4cd-434d-a958-884d6980c459">
	//         <cybox:Object id="ExampleNamespace:Object-1980ce43-8e03-490b-863a-ea404d12242e">
	//           <cybox:Properties xsi:type="AddressObject:AddressObjectType">
	//             <AddressObject:Address_Value condition="Equals" apply_condition="ANY">10.0.0.0##comma##10.0.0.1##comma##10.0.0.2</AddressObject:Address_Value>
	//           </cybox:Properties>
	//         </cybox:Object>
	//       </indicator:Observable>
	//     </stix:Indicator>
	//   </stix:Indicators>
	// </stix:STIX_Package>
}
