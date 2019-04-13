$(document).ready(function () {

    // On page load - if additional sections exist
    // dynamically add IDs to the selector and textarea
    $(function() {
        var additionalSections = $('#additionalSections').children();

        // Loop each additional section
        $.each(additionalSections, function(index, item) {
            // Increment ID
            newID = index+1;

            // Set the selectors ID and name attributes
            let select = $(item).find("select");
            select.id = "additionalCVSectionSubject" + newID;
            select.attr({
                "name": "additionalCVSectionSubject" + newID

            });

            // Set the textareas ID and name attributes
            let textarea = $(item).find("textarea");
            textarea.id = "additionalCVSectionValue" + newID;
            textarea.attr("name", "additionalCVSectionValue" + newID);

            // Set the ID for the section
            item.id = newID


        });
    });

    // Function to add new CV section
    $('#addSection').on('click', function () {
        var section_count = $('#additionalSections').children().length + 1;
        var new_id = parseInt($('#additionalSections').children().last().attr('id'));

        // If the id cannot be found, then it is the first additional section added
        // Assign 1, otherwise add 1
        if (isNaN(new_id)) {
            new_id = 1;
        } else {
            new_id++;
        }

        // Create the elements to insert into the DOM
        var card = document.createElement('div');
        card.className = "card mb-2";
        card.id = new_id;

        var card_body = document.createElement('div');
        card_body.className = "card-body";

        var remove_button = document.createElement('button');
        remove_button.className = "btn btn-danger";
        remove_button.type = "button";
        remove_button.id = "removeSection";
        remove_button.innerHTML = "<span class=\"fa fa-trash\" aria-hidden=\"true\"></span> Remove";

        var data_type_form_group = document.createElement('div');
        data_type_form_group.className = "form-group";

        var select_label = document.createElement('label');
        select_label.setAttribute("for", "additionalCVSectionSubject" + new_id);
        select_label.innerHTML = "Subject";

        var select = document.createElement("select");
        select.className = "form-control";
        select.id = "additionalCVSectionSubject" + new_id;
        select.setAttribute("name", "additionalCVSectionSubject" + new_id);

        // List of the values in the select dropdown. To add new values, update this
        var selectValues = {
            "Experience": "Experience",
            "Education": "Education",
            "Skills": "Skills",
            "Certifications": "Certifications",
            "Traits": "Traits",
            "Interests": "Interests",
            "Other": "Other"
        };

        // Function to loop through each value and create the new option elements
        $.each(selectValues, function (key, value) {
            var option = document.createElement('option');
            option.innerHTML = value;
            select.append(option)
        });

        var data_value_form_group = document.createElement('div');
        data_value_form_group.className = "form-group";

        var textarea_label = document.createElement('label');
        textarea_label.setAttribute("for", "additionalCVSectionValue" + new_id);
        textarea_label.innerHTML = "Details";

        var textarea = document.createElement("textarea");
        textarea.className = "form-control";
        textarea.id = "additionalCVSectionValue" + new_id;
        textarea.setAttribute("rows", "3");
        textarea.setAttribute("name", "additionalCVSectionValue" + new_id);

        // data_type_form_group contains the selector dropdown
        data_type_form_group.append(select_label);
        data_type_form_group.append(select);

        // data_value_form_group contains the textarea
        data_value_form_group.append(textarea_label);
        data_value_form_group.append(textarea);
        card_body.append(data_type_form_group);
        card_body.append(data_value_form_group);
        card_body.append(remove_button);
        card.append(card_body);

        // If there is more than 5 sections, prevent any more from being added
        if (section_count < 7) {
            $('#additionalSections').append(card);
        }
    });

    // Function to remove selected CV section
    $('#additionalSections').on('click', '#removeSection', function () {
        $(this).closest('.card').remove();
    });
});
